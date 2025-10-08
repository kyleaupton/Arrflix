package service

import (
	"context"

	dbgen "github.com/kyleaupton/snaggle/backend/internal/db/sqlc"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
)

type MediaService struct {
	repo *repo.Repository
}

func NewMediaService(r *repo.Repository) *MediaService {
	return &MediaService{repo: r}
}

func (s *MediaService) List(ctx context.Context) ([]dbgen.MediaItem, error) {
	return s.repo.ListMediaItems(ctx)
}
