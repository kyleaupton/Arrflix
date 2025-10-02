package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/snaggle/backend/internal/db/sqlc"
)

type AuthRepo interface {
	GetUserByEmail(ctx context.Context, email string) (dbgen.AppUser, error)
	UpdateUserPassword(ctx context.Context, userID pgtype.UUID, newHash string) error
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (dbgen.AppUser, error) {
	return r.Q.GetUserByEmail(ctx, email)
}

func (r *Repository) UpdateUserPassword(ctx context.Context, userID pgtype.UUID, newHash string) error {
	return r.Q.UpdateUserPassword(ctx, dbgen.UpdateUserPasswordParams{ID: userID, PasswordHash: &newHash})
}
