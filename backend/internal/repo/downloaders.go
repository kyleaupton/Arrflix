package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/Arrflix/internal/db/sqlc"
)

type DownloaderRepo interface {
	ListDownloaders(ctx context.Context) ([]dbgen.Downloader, error)
	GetDownloader(ctx context.Context, id pgtype.UUID) (dbgen.Downloader, error)
	GetDefaultDownloader(ctx context.Context, protocol string) (dbgen.Downloader, error)
	CreateDownloader(ctx context.Context, name, downloaderType, protocol, url string, username, password *string, configJSON []byte, enabled, isDefault bool) (dbgen.Downloader, error)
	UpdateDownloader(ctx context.Context, id pgtype.UUID, name, downloaderType, protocol, url string, username, password *string, configJSON []byte, enabled, isDefault bool) (dbgen.Downloader, error)
	DeleteDownloader(ctx context.Context, id pgtype.UUID) error
}

func (r *Repository) ListDownloaders(ctx context.Context) ([]dbgen.Downloader, error) {
	return r.Q.ListDownloaders(ctx)
}

func (r *Repository) GetDownloader(ctx context.Context, id pgtype.UUID) (dbgen.Downloader, error) {
	return r.Q.GetDownloader(ctx, id)
}

func (r *Repository) GetDefaultDownloader(ctx context.Context, protocol string) (dbgen.Downloader, error) {
	return r.Q.GetDefaultDownloader(ctx, protocol)
}

func (r *Repository) CreateDownloader(ctx context.Context, name, downloaderType, protocol, url string, username, password *string, configJSON []byte, enabled, isDefault bool) (dbgen.Downloader, error) {
	var configJSONVal []byte
	if configJSON != nil && len(configJSON) > 0 {
		configJSONVal = configJSON
	}

	return r.Q.CreateDownloader(ctx, dbgen.CreateDownloaderParams{
		Name:           name,
		DownloaderType: downloaderType,
		Protocol:       protocol,
		Url:            url,
		Username:       username,
		Password:       password,
		ConfigJson:     configJSONVal,
		Enabled:        enabled,
		IsDefault:      isDefault,
	})
}

func (r *Repository) UpdateDownloader(ctx context.Context, id pgtype.UUID, name, downloaderType, protocol, url string, username, password *string, configJSON []byte, enabled, isDefault bool) (dbgen.Downloader, error) {
	var configJSONVal []byte
	if configJSON != nil && len(configJSON) > 0 {
		configJSONVal = configJSON
	}

	return r.Q.UpdateDownloader(ctx, dbgen.UpdateDownloaderParams{
		ID:             id,
		Name:           name,
		DownloaderType: downloaderType,
		Protocol:       protocol,
		Url:            url,
		Username:       username,
		Password:       password,
		ConfigJson:     configJSONVal,
		Enabled:        enabled,
		IsDefault:      isDefault,
	})
}

func (r *Repository) DeleteDownloader(ctx context.Context, id pgtype.UUID) error {
	return r.Q.DeleteDownloader(ctx, id)
}
