package reconciler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kyleaupton/snaggle/ops/internal/config"
	"github.com/kyleaupton/snaggle/ops/internal/docker"
	"github.com/kyleaupton/snaggle/ops/internal/services"
)

// Reconciler manages the reconciliation loop for all services
type Reconciler struct {
	docker *docker.Client
	config *config.Config
	static []services.Service
	db     *pgxpool.Pool
}

// New creates a new reconciler
func New(dockerClient *docker.Client, cfg *config.Config) *Reconciler {
	// Initialize static services
	static := []services.Service{
		services.NewPostgres(cfg),
		services.NewProwlarr(cfg),
		services.NewAPI(cfg),
		services.NewNginx(cfg),
	}

	return &Reconciler{
		docker: dockerClient,
		config: cfg,
		static: static,
	}
}

// Run starts the reconciliation loop
func (r *Reconciler) Run(ctx context.Context) error {
	// Connect to database
	db, err := pgxpool.New(ctx, r.config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()
	r.db = db

	log.Printf("Starting reconciler with interval: %v", r.config.ReconcileInterval)

	ticker := time.NewTicker(r.config.ReconcileInterval)
	defer ticker.Stop()

	// Run initial reconciliation
	if err := r.Reconcile(ctx); err != nil {
		log.Printf("Initial reconciliation failed: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Reconciler stopping...")
			return ctx.Err()
		case <-ticker.C:
			if err := r.Reconcile(ctx); err != nil {
				log.Printf("Reconciliation failed: %v", err)
			}
		}
	}
}

// Reconcile performs a single reconciliation cycle
func (r *Reconciler) Reconcile(ctx context.Context) error {
	log.Println("Starting reconciliation cycle")

	// 1. Ensure network and volumes exist
	if err := r.docker.EnsureNetwork(ctx, r.config.NetworkName); err != nil {
		return fmt.Errorf("failed to ensure network: %w", err)
	}

	// Ensure required volumes
	volumes := []string{"snaggle_pg_data", "snaggle_prowlarr_data"}
	for _, vol := range volumes {
		if err := r.docker.EnsureVolume(ctx, vol); err != nil {
			return fmt.Errorf("failed to ensure volume %s: %w", vol, err)
		}
	}

	// 2. Build desired state
	desiredServices, err := r.buildDesiredState(ctx)
	if err != nil {
		return fmt.Errorf("failed to build desired state: %w", err)
	}

	// 3. Get actual state
	actualContainers, err := r.docker.ListManagedContainers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get actual state: %w", err)
	}

	// 4. Build dependency graph and reconcile
	return r.reconcileServices(ctx, desiredServices, actualContainers)
}

// buildDesiredState builds the desired state from static services and DB
func (r *Reconciler) buildDesiredState(ctx context.Context) ([]services.Service, error) {
	var desired []services.Service

	// Add static services
	desired = append(desired, r.static...)

	// Add dynamic services from database
	dynamicServices, err := r.getDynamicServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get dynamic services: %w", err)
	}
	desired = append(desired, dynamicServices...)

	return desired, nil
}

// getDynamicServices queries the database for enabled service instances
func (r *Reconciler) getDynamicServices(ctx context.Context) ([]services.Service, error) {
	query := `SELECT id, name, type, enabled, config FROM service_instance WHERE enabled = true`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query service instances: %w", err)
	}
	defer rows.Close()

	var serviceList []services.Service
	for rows.Next() {
		var instance services.ServiceInstance
		var configJSON []byte

		err := rows.Scan(&instance.ID, &instance.Name, &instance.Type, &instance.Enabled, &configJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service instance: %w", err)
		}

		// Parse JSON config
		instance.Config = make(map[string]interface{})
		// TODO: Parse JSON config properly
		// For now, use empty config

		// Create dynamic service
		service, err := services.NewDynamicService(r.config, &instance)
		if err != nil {
			log.Printf("Failed to create dynamic service %s: %v", instance.Name, err)
			continue
		}

		serviceList = append(serviceList, service)
	}

	return serviceList, nil
}

// reconcileServices reconciles the desired state with actual state
func (r *Reconciler) reconcileServices(ctx context.Context, desired []services.Service, actual []*docker.ContainerStatus) error {
	// Build maps for easier lookup
	desiredMap := make(map[string]services.Service)
	for _, service := range desired {
		desiredMap[service.Name()] = service
	}

	actualMap := make(map[string]*docker.ContainerStatus)
	for _, container := range actual {
		actualMap[container.Name] = container
	}

	// Build dependency graph and process in order
	dependencyOrder, err := r.buildDependencyOrder(desired)
	if err != nil {
		return fmt.Errorf("failed to build dependency order: %w", err)
	}

	// Process services in dependency order
	for _, serviceName := range dependencyOrder {
		service := desiredMap[serviceName]
		if err := r.reconcileService(ctx, service, actualMap[serviceName]); err != nil {
			log.Printf("Failed to reconcile service %s: %v", serviceName, err)
		}
	}

	// Remove containers that are no longer desired
	for _, container := range actual {
		if _, exists := desiredMap[container.Name]; !exists {
			log.Printf("Removing unwanted container: %s", container.Name)
			if err := r.docker.StopContainer(ctx, container.Name); err != nil {
				log.Printf("Failed to stop container %s: %v", container.Name, err)
			}
			if err := r.docker.RemoveContainer(ctx, container.Name); err != nil {
				log.Printf("Failed to remove container %s: %v", container.Name, err)
			}
		}
	}

	return nil
}

// reconcileService reconciles a single service
func (r *Reconciler) reconcileService(ctx context.Context, service services.Service, actual *docker.ContainerStatus) error {
	// Check dependencies first
	for _, dep := range service.DependsOn() {
		depStatus, err := r.docker.GetContainerStatus(ctx, dep)
		if err != nil {
			log.Printf("Dependency %s not ready for %s: %v", dep, service.Name(), err)
			return fmt.Errorf("dependency %s not ready", dep)
		}
		if depStatus.Status != "running" {
			log.Printf("Dependency %s not running for %s", dep, service.Name())
			return fmt.Errorf("dependency %s not running", dep)
		}
	}

	// Ensure container exists and is running
	if actual == nil {
		log.Printf("Creating container: %s", service.Name())
		return r.docker.EnsureContainer(ctx, service)
	}

	if actual.Status != "running" {
		log.Printf("Starting container: %s", service.Name())
		return r.docker.EnsureContainer(ctx, service)
	}

	log.Printf("Container %s is running", service.Name())
	return nil
}

// buildDependencyOrder builds a topological sort of services based on dependencies
func (r *Reconciler) buildDependencyOrder(services []services.Service) ([]string, error) {
	// Simple implementation - just return services in order
	// TODO: Implement proper topological sort
	var order []string
	for _, service := range services {
		order = append(order, service.Name())
	}
	return order, nil
}
