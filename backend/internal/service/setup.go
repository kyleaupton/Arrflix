package service

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	dbgen "github.com/kyleaupton/arrflix/internal/db/sqlc"
	pw "github.com/kyleaupton/arrflix/internal/password"
	"github.com/kyleaupton/arrflix/internal/repo"
)

var (
	ErrAlreadyInitialized = errors.New("system already initialized")
	ErrSetupFailed        = errors.New("setup failed")
)

type SetupService struct {
	repo  *repo.Repository
	users *UsersService
}

func NewSetupService(r *repo.Repository, users *UsersService) *SetupService {
	return &SetupService{repo: r, users: users}
}

// IsInitialized checks if the system has been initialized
func (s *SetupService) IsInitialized(ctx context.Context) (bool, error) {
	return s.repo.GetSystemInitialized(ctx)
}

// Initialize performs the one-time setup operation atomically:
// 1. Check initialization status
// 2. Create admin user
// 3. Assign admin role
// 4. Mark system as initialized
// All operations occur in a single transaction for atomicity.
func (s *SetupService) Initialize(ctx context.Context, email, username, password string) error {
	// Start transaction for atomicity
	tx, err := s.repo.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable, // Highest isolation to prevent concurrent setup
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Check if already initialized (within transaction)
	txQueries := s.repo.Q.WithTx(tx)
	initialized, err := txQueries.GetSystemInitialized(ctx)
	if err != nil {
		return err
	}
	if initialized {
		return ErrAlreadyInitialized
	}

	// Validate input
	email = strings.TrimSpace(email)
	username = strings.TrimSpace(username)
	if email == "" {
		return errors.New("email required")
	}
	if username == "" {
		return errors.New("username required")
	}
	if password == "" {
		return errors.New("password required")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	// Hash password
	passwordHash, err := pw.Hash(password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Create admin user within transaction
	user, err := txQueries.CreateUser(ctx, dbgen.CreateUserParams{
		Email:        &email,
		Username:     username,
		PasswordHash: &passwordHash,
		IsActive:     true,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return errors.New("email already exists")
		}
		return err
	}

	// Assign admin role within transaction
	role, err := txQueries.GetRoleByName(ctx, "admin")
	if err != nil {
		return errors.New("admin role not found")
	}
	if err := txQueries.AssignRole(ctx, dbgen.AssignRoleParams{
		UserID: user.ID,
		RoleID: role.ID,
	}); err != nil {
		return err
	}

	// Mark system as initialized (conditional update prevents double-init)
	if err := txQueries.SetSystemInitialized(ctx); err != nil {
		return err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return ErrSetupFailed
	}

	return nil
}
