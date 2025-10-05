package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/snaggle/backend/internal/db/sqlc"
)

type LibraryRepo interface {
	ListLibraries(ctx context.Context) ([]dbgen.Library, error)
	GetLibrary(ctx context.Context, id pgtype.UUID) (dbgen.Library, error)
	CreateLibrary(ctx context.Context, name, typ, rootPath string, enabled *bool) (dbgen.Library, error)
	UpdateLibrary(ctx context.Context, id pgtype.UUID, name, typ, rootPath string, enabled bool) (dbgen.Library, error)
	DeleteLibrary(ctx context.Context, id pgtype.UUID) error
}

func (r *Repository) ListLibraries(ctx context.Context) ([]dbgen.Library, error) {
	return r.Q.ListLibraries(ctx)
}

func (r *Repository) GetLibrary(ctx context.Context, id pgtype.UUID) (dbgen.Library, error) {
	return r.Q.GetLibrary(ctx, id)
}

func (r *Repository) CreateLibrary(ctx context.Context, name, typ, rootPath string, enabled bool) (dbgen.Library, error) {
	return r.Q.CreateLibrary(ctx, dbgen.CreateLibraryParams{
		Name:     name,
		Type:     typ,
		RootPath: rootPath,
		Enabled:  enabled,
	})
}

func (r *Repository) UpdateLibrary(ctx context.Context, id pgtype.UUID, name, typ, rootPath string, enabled bool) (dbgen.Library, error) {
	return r.Q.UpdateLibrary(ctx, dbgen.UpdateLibraryParams{
		ID:       id,
		Name:     name,
		Type:     typ,
		RootPath: rootPath,
		Enabled:  enabled,
	})
}

func (r *Repository) DeleteLibrary(ctx context.Context, id pgtype.UUID) error {
	return r.Q.DeleteLibrary(ctx, id)
}
