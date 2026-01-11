package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/Arrflix/internal/db/sqlc"
)

type AuthRepo interface {
	GetUserByEmail(ctx context.Context, email string) (dbgen.AppUser, error)
	UpdateUserPassword(ctx context.Context, userID pgtype.UUID, newHash string) error
	// User CRUD
	ListUsers(ctx context.Context) ([]dbgen.ListUsersRow, error)
	GetUser(ctx context.Context, id pgtype.UUID) (dbgen.GetUserRow, error)
	CreateUser(ctx context.Context, email, displayName, passwordHash string, isActive bool) (dbgen.AppUser, error)
	UpdateUser(ctx context.Context, id pgtype.UUID, email, displayName string, isActive bool) (dbgen.AppUser, error)
	DeleteUser(ctx context.Context, id pgtype.UUID) error
	// Role Management
	ListRoles(ctx context.Context) ([]dbgen.Role, error)
	ListUserRoles(ctx context.Context, userID pgtype.UUID) ([]dbgen.Role, error)
	GetRoleByName(ctx context.Context, name string) (dbgen.Role, error)
	AssignRole(ctx context.Context, userID, roleID pgtype.UUID) error
	UnassignAllRoles(ctx context.Context, userID pgtype.UUID) error
	CountUsersByRole(ctx context.Context, roleID pgtype.UUID) (int64, error)
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (dbgen.AppUser, error) {
	return r.Q.GetUserByEmail(ctx, email)
}

func (r *Repository) UpdateUserPassword(ctx context.Context, userID pgtype.UUID, newHash string) error {
	return r.Q.UpdateUserPassword(ctx, dbgen.UpdateUserPasswordParams{ID: userID, PasswordHash: &newHash})
}

// User CRUD implementations

func (r *Repository) ListUsers(ctx context.Context) ([]dbgen.ListUsersRow, error) {
	return r.Q.ListUsers(ctx)
}

func (r *Repository) GetUser(ctx context.Context, id pgtype.UUID) (dbgen.GetUserRow, error) {
	return r.Q.GetUser(ctx, id)
}

func (r *Repository) CreateUser(ctx context.Context, email, displayName, passwordHash string, isActive bool) (dbgen.AppUser, error) {
	return r.Q.CreateUser(ctx, dbgen.CreateUserParams{
		Email:        &email,
		DisplayName:  &displayName,
		PasswordHash: &passwordHash,
		IsActive:     isActive,
	})
}

func (r *Repository) UpdateUser(ctx context.Context, id pgtype.UUID, email, displayName string, isActive bool) (dbgen.AppUser, error) {
	return r.Q.UpdateUser(ctx, dbgen.UpdateUserParams{
		ID:          id,
		Email:       &email,
		DisplayName: &displayName,
		IsActive:    isActive,
	})
}

func (r *Repository) DeleteUser(ctx context.Context, id pgtype.UUID) error {
	return r.Q.DeleteUser(ctx, id)
}

// Role Management implementations

func (r *Repository) ListRoles(ctx context.Context) ([]dbgen.Role, error) {
	return r.Q.ListRoles(ctx)
}

func (r *Repository) ListUserRoles(ctx context.Context, userID pgtype.UUID) ([]dbgen.Role, error) {
	return r.Q.ListUserRoles(ctx, userID)
}

func (r *Repository) GetRoleByName(ctx context.Context, name string) (dbgen.Role, error) {
	return r.Q.GetRoleByName(ctx, name)
}

func (r *Repository) AssignRole(ctx context.Context, userID, roleID pgtype.UUID) error {
	return r.Q.AssignRole(ctx, dbgen.AssignRoleParams{
		UserID: userID,
		RoleID: roleID,
	})
}

func (r *Repository) UnassignAllRoles(ctx context.Context, userID pgtype.UUID) error {
	return r.Q.UnassignAllRoles(ctx, userID)
}

func (r *Repository) CountUsersByRole(ctx context.Context, roleID pgtype.UUID) (int64, error) {
	return r.Q.CountUsersByRole(ctx, roleID)
}
