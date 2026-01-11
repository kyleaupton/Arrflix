package service

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/Arrflix/internal/db/sqlc"
	pw "github.com/kyleaupton/Arrflix/internal/password"
	"github.com/kyleaupton/Arrflix/internal/repo"
)

type UsersService struct {
	repo *repo.Repository
}

func NewUsersService(r *repo.Repository) *UsersService {
	return &UsersService{repo: r}
}

// List returns all users with their roles
func (s *UsersService) List(ctx context.Context) ([]dbgen.ListUsersRow, error) {
	return s.repo.ListUsers(ctx)
}

// Get returns a single user with roles
func (s *UsersService) Get(ctx context.Context, id pgtype.UUID) (dbgen.GetUserRow, error) {
	return s.repo.GetUser(ctx, id)
}

// Create creates a new user with password and role assignment
func (s *UsersService) Create(ctx context.Context, email, displayName, password string, roleName string, isActive bool) (dbgen.AppUser, error) {
	// Validation
	email = strings.TrimSpace(email)
	displayName = strings.TrimSpace(displayName)

	if email == "" {
		return dbgen.AppUser{}, errors.New("email required")
	}
	if displayName == "" {
		return dbgen.AppUser{}, errors.New("display_name required")
	}
	if password == "" {
		return dbgen.AppUser{}, errors.New("password required")
	}
	if len(password) < 8 {
		return dbgen.AppUser{}, errors.New("password must be at least 8 characters")
	}
	if roleName == "" {
		roleName = "user" // Default to 'user' role
	}

	// Hash password
	passwordHash, err := pw.Hash(password)
	if err != nil {
		return dbgen.AppUser{}, errors.New("failed to hash password")
	}

	// Create user
	user, err := s.repo.CreateUser(ctx, email, displayName, passwordHash, isActive)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return dbgen.AppUser{}, errors.New("email already exists")
		}
		return dbgen.AppUser{}, err
	}

	// Assign role
	role, err := s.repo.GetRoleByName(ctx, roleName)
	if err != nil {
		// User created but role assignment failed - still return user
		return user, nil
	}

	_ = s.repo.AssignRole(ctx, user.ID, role.ID)

	return user, nil
}

// Update updates user fields (not password)
func (s *UsersService) Update(ctx context.Context, id pgtype.UUID, email, displayName string, isActive bool) (dbgen.AppUser, error) {
	email = strings.TrimSpace(email)
	displayName = strings.TrimSpace(displayName)

	if email == "" {
		return dbgen.AppUser{}, errors.New("email required")
	}
	if displayName == "" {
		return dbgen.AppUser{}, errors.New("display_name required")
	}

	return s.repo.UpdateUser(ctx, id, email, displayName, isActive)
}

// UpdatePassword changes a user's password
func (s *UsersService) UpdatePassword(ctx context.Context, userID pgtype.UUID, newPassword string) error {
	if newPassword == "" {
		return errors.New("password required")
	}
	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	hash, err := pw.Hash(newPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	return s.repo.UpdateUserPassword(ctx, userID, hash)
}

// AssignRole assigns a single role to a user (replaces existing roles)
func (s *UsersService) AssignRole(ctx context.Context, userID pgtype.UUID, roleName string) error {
	role, err := s.repo.GetRoleByName(ctx, roleName)
	if err != nil {
		return errors.New("role not found")
	}

	// Unassign all existing roles first
	if err := s.repo.UnassignAllRoles(ctx, userID); err != nil {
		return err
	}

	return s.repo.AssignRole(ctx, userID, role.ID)
}

// Delete removes a user
func (s *UsersService) Delete(ctx context.Context, id pgtype.UUID) error {
	// Check if this is the last admin
	roles, err := s.repo.ListUserRoles(ctx, id)
	if err == nil {
		for _, role := range roles {
			if role.Name == "admin" {
				// Count admins
				adminRole, _ := s.repo.GetRoleByName(ctx, "admin")
				count, _ := s.repo.CountUsersByRole(ctx, adminRole.ID)
				if count <= 1 {
					return errors.New("cannot delete the last admin user")
				}
			}
		}
	}

	return s.repo.DeleteUser(ctx, id)
}

// ListRoles returns all available roles
func (s *UsersService) ListRoles(ctx context.Context) ([]dbgen.Role, error) {
	return s.repo.ListRoles(ctx)
}
