package repo

import (
	"context"

	dbgen "github.com/kyleaupton/snaggle/backend/internal/db/sqlc"
)

type AuthRepo interface {
    GetUserByEmail(ctx context.Context, email string) (dbgen.AppUser, error)
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (dbgen.AppUser, error) {
    return r.Q.GetUserByEmail(ctx, email)
}

