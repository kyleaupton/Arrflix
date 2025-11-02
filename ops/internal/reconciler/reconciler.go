package reconciler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	// Initialize static services based on runtime mode
	var static []services.Service

	if cfg.RuntimeMode == "dev" {
		static = []services.Service{
			services.NewPostgres(cfg),
			services.NewProwlarr(cfg),
			services.NewAPIDev(cfg),
			services.NewWebDev(cfg),
			services.NewNginxDev(cfg),
		}
	} else {
		static = []services.Service{
			services.NewPostgres(cfg),
			services.NewProwlarr(cfg),
			services.NewAPI(cfg),
			services.NewNginx(cfg),
		}
	}

	return &Reconciler{
		docker: dockerClient,
		config: cfg,
		static: static,
	}
}

// Run starts the reconciliation loop
func (r *Reconciler) Run(ctx context.Context) error {
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
			if r.db != nil {
				r.db.Close()
			}
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
	log.Printf("Network %s ensured", r.config.NetworkName)

	// Ensure required volumes
	volumes := []string{"snaggle_pg_data", "snaggle_prowlarr_data"}
	for _, vol := range volumes {
		if err := r.docker.EnsureVolume(ctx, vol); err != nil {
			return fmt.Errorf("failed to ensure volume %s: %w", vol, err)
		}
	}
	log.Printf("Volumes ensured: %v", volumes)

	// 2. Build desired state
	desiredServices, err := r.buildDesiredState(ctx)
	if err != nil {
		return fmt.Errorf("failed to build desired state: %w", err)
	}
	log.Printf("Desired state built with %d services", len(desiredServices))

	// 3. Get actual state
	actualContainers, err := r.docker.ListManagedContainers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get actual state: %w", err)
	}
	log.Printf("Found %d actual containers", len(actualContainers))

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
	// Try to connect to database if not already connected
	if r.db == nil {
		db, err := pgxpool.New(ctx, r.config.DatabaseURL)
		if err != nil {
			// Database not available yet, return empty list
			log.Printf("Database not available yet, skipping dynamic services: %v", err)
			return []services.Service{}, nil
		}
		r.db = db
	}

	query := `SELECT id, name, type, enabled, config FROM service_instance WHERE enabled = true`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		// Database query failed, return empty list
		log.Printf("Failed to query service instances, skipping dynamic services: %v", err)
		return []services.Service{}, nil
	}
	defer rows.Close()

	var serviceList []services.Service
	for rows.Next() {
		var instance services.ServiceInstance
		var configJSON []byte

		err := rows.Scan(&instance.ID, &instance.Name, &instance.Type, &instance.Enabled, &configJSON)
		if err != nil {
			log.Printf("Failed to scan service instance: %v", err)
			continue
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
		// Wait for dependency to be healthy (this handles all waiting logic)
		if err := r.waitForDependencyHealthy(ctx, dep); err != nil {
			log.Printf("Dependency %s not healthy for %s: %v", dep, service.Name(), err)
			return fmt.Errorf("dependency %s not healthy", dep)
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

// waitForDependencyHealthy waits for a dependency container to be healthy
func (r *Reconciler) waitForDependencyHealthy(ctx context.Context, depName string) error {
	log.Printf("Waiting for dependency %s to be healthy...", depName)

	// Poll every second for up to 60 seconds
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeout := time.After(60 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return fmt.Errorf("timeout waiting for dependency %s to be healthy", depName)
		case <-ticker.C:
			// First check if container exists and is running
			status, err := r.docker.GetContainerStatus(ctx, depName)
			if err != nil {
				log.Printf("Dependency %s container not found yet: %v", depName, err)
				continue
			}
			// Check if container is running (Docker status can be "Up X seconds (healthy)" etc.)
			if !strings.Contains(status.Status, "Up") && status.Status != "running" {
				log.Printf("Dependency %s not running yet (status: %s)", depName, status.Status)
				continue
			}

			// Container is running, now check if it's healthy
			healthy, err := r.checkDependencyHealth(ctx, depName)
			if err != nil {
				log.Printf("Health check failed for %s: %v", depName, err)
				continue
			}
			if healthy {
				log.Printf("Dependency %s is now healthy", depName)
				return nil
			}
			log.Printf("Dependency %s not yet healthy, retrying...", depName)
		}
	}
}

// checkDependencyHealth performs a health check on a specific dependency
func (r *Reconciler) checkDependencyHealth(ctx context.Context, depName string) (bool, error) {
	switch depName {
	case "snaggle-postgres":
		return r.checkPostgresHealth(ctx)
	case "snaggle-prowlarr":
		return r.checkProwlarrHealth(ctx)
	case "snaggle-api-dev":
		return r.checkAPIHealth(ctx)
	default:
		// For unknown dependencies, just check if container is running
		status, err := r.docker.GetContainerStatus(ctx, depName)
		if err != nil {
			return false, err
		}
		return status.Status == "running", nil
	}
}

// checkPostgresHealth checks if PostgreSQL is ready to accept connections
func (r *Reconciler) checkPostgresHealth(ctx context.Context) (bool, error) {
	// Try to connect to PostgreSQL
	db, err := pgxpool.New(ctx, r.config.DatabaseURL)
	if err != nil {
		return false, nil // Not ready yet, but not an error
	}
	defer db.Close()

	// Test the connection
	var result int
	err = db.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return false, nil // Not ready yet, but not an error
	}

	return result == 1, nil
}

// checkProwlarrHealth checks if Prowlarr is responding
func (r *Reconciler) checkProwlarrHealth(ctx context.Context) (bool, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Create request to Prowlarr API
	req, err := http.NewRequestWithContext(ctx, "GET", "http://snaggle-prowlarr:9696/api", nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	// Add API key header
	req.Header.Set("X-Api-Key", r.config.ProwlarrAPIKey)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return false, nil // Not ready yet, but not a critical error
	}
	defer resp.Body.Close()

	// Check if response is successful (200 OK)
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	// Not ready yet
	return false, nil
}

// checkAPIHealth checks if the API is responding
func (r *Reconciler) checkAPIHealth(ctx context.Context) (bool, error) {
	// For now, just check if container is running
	// TODO: Add actual HTTP health check to API
	status, err := r.docker.GetContainerStatus(ctx, "snaggle-api-dev")
	if err != nil {
		return false, err
	}
	return status.Status == "running", nil
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
