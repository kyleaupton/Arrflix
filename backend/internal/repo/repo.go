package repo

import (
	"github.com/jackc/pgx/v5/pgxpool"
	dbgen "github.com/kyleaupton/snaggle/backend/internal/db/sqlc"
)

type Repository struct {
    Pool *pgxpool.Pool
    Q    *dbgen.Queries
}

func New(pool *pgxpool.Pool) *Repository {
    return &Repository{
        Pool: pool,
        Q:    dbgen.New(pool),
    }
}
