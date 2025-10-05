package service

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/snaggle/backend/internal/db/sqlc"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
)

type LibrariesService struct {
	repo *repo.Repository
}

func NewLibrariesService(r *repo.Repository) *LibrariesService {
	return &LibrariesService{repo: r}
}

func (s *LibrariesService) List(ctx context.Context) ([]dbgen.Library, error) {
	return s.repo.ListLibraries(ctx)
}

func (s *LibrariesService) Get(ctx context.Context, id pgtype.UUID) (dbgen.Library, error) {
	return s.repo.GetLibrary(ctx, id)
}

func (s *LibrariesService) Create(ctx context.Context, name, typ, rootPath string, enabled bool) (dbgen.Library, error) {
	if name == "" {
		return dbgen.Library{}, errors.New("name required")
	}
	if typ != "movie" && typ != "series" {
		return dbgen.Library{}, errors.New("type must be 'movie' or 'series'")
	}
	if rootPath == "" {
		return dbgen.Library{}, errors.New("root_path required")
	}
	if _, err := os.Stat(rootPath); err != nil {
		return dbgen.Library{}, errors.New("root_path not found on server")
	}
	return s.repo.CreateLibrary(ctx, name, typ, rootPath, enabled)
}

func (s *LibrariesService) Update(ctx context.Context, id pgtype.UUID, name, typ, rootPath string, enabled bool) (dbgen.Library, error) {
	if name == "" {
		return dbgen.Library{}, errors.New("name required")
	}
	if typ != "movie" && typ != "series" {
		return dbgen.Library{}, errors.New("type must be 'movie' or 'series'")
	}
	if rootPath == "" {
		return dbgen.Library{}, errors.New("root_path required")
	}
	if _, err := os.Stat(rootPath); err != nil {
		return dbgen.Library{}, errors.New("root_path not found on server")
	}
	return s.repo.UpdateLibrary(ctx, id, name, typ, rootPath, enabled)
}

func (s *LibrariesService) Delete(ctx context.Context, id pgtype.UUID) error {
	return s.repo.DeleteLibrary(ctx, id)
}
