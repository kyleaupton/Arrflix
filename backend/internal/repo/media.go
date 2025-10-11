package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/snaggle/backend/internal/db/sqlc"
)

type MediaRepo interface {
	// Media items
	ListMediaItems(ctx context.Context) ([]dbgen.MediaItem, error)
	GetMediaItem(ctx context.Context, id pgtype.UUID) (dbgen.MediaItem, error)
	GetMediaItemByTmdbID(ctx context.Context, tmdbID int64) (dbgen.MediaItem, error)
	CreateMediaItem(ctx context.Context, libraryID pgtype.UUID, typ, title string, year *int32, tmdbID *int32) (dbgen.MediaItem, error)
	UpdateMediaItem(ctx context.Context, id pgtype.UUID, title string, year *int32, tmdbID *int32) (dbgen.MediaItem, error)
	DeleteMediaItem(ctx context.Context, id pgtype.UUID) error

	// Seasons
	ListSeasonsForMedia(ctx context.Context, mediaItemID pgtype.UUID) ([]dbgen.MediaSeason, error)
	UpsertSeason(ctx context.Context, mediaItemID pgtype.UUID, seasonNumber int32, airDate pgtype.Date) (dbgen.MediaSeason, error)

	// Episodes
	ListEpisodesForSeason(ctx context.Context, seasonID pgtype.UUID) ([]dbgen.MediaEpisode, error)
	UpsertEpisode(ctx context.Context, seasonID pgtype.UUID, episodeNumber int32, title *string, airDate pgtype.Date, tmdbID *int64, tvdbID *int64) (dbgen.MediaEpisode, error)

	// Files
	GetMediaFileByPath(ctx context.Context, path string) (dbgen.MediaFile, error)
	CreateMediaFile(ctx context.Context, mediaItemID pgtype.UUID, seasonID, episodeID *pgtype.UUID, path string) (dbgen.MediaFile, error)
	DeleteMediaFile(ctx context.Context, id pgtype.UUID) error
}

func (r *Repository) ListMediaItems(ctx context.Context) ([]dbgen.MediaItem, error) {
	return r.Q.ListMediaItems(ctx)
}

func (r *Repository) GetMediaItem(ctx context.Context, id pgtype.UUID) (dbgen.MediaItem, error) {
	return r.Q.GetMediaItem(ctx, id)
}

func (r *Repository) GetMediaItemByTmdbID(ctx context.Context, tmdbID int64) (dbgen.MediaItem, error) {
	return r.Q.GetMediaItemByTmdbID(ctx, &tmdbID)
}

func (r *Repository) CreateMediaItem(ctx context.Context, libraryID pgtype.UUID, typ, title string, year *int32, tmdbID *int64) (dbgen.MediaItem, error) {
	return r.Q.CreateMediaItem(ctx, dbgen.CreateMediaItemParams{
		LibraryID: libraryID,
		Type:      typ,
		Title:     title,
		Year:      year,
		TmdbID:    tmdbID,
	})
}

func (r *Repository) UpdateMediaItem(ctx context.Context, id pgtype.UUID, title string, year *int32, tmdbID *int64) (dbgen.MediaItem, error) {
	return r.Q.UpdateMediaItem(ctx, dbgen.UpdateMediaItemParams{
		ID:     id,
		Title:  title,
		Year:   year,
		TmdbID: tmdbID,
	})
}

func (r *Repository) DeleteMediaItem(ctx context.Context, id pgtype.UUID) error {
	return r.Q.DeleteMediaItem(ctx, id)
}

func (r *Repository) ListSeasonsForMedia(ctx context.Context, mediaID pgtype.UUID) ([]dbgen.MediaSeason, error) {
	return r.Q.ListSeasonsForMedia(ctx, mediaID)
}

func (r *Repository) UpsertSeason(ctx context.Context, mediaItemID pgtype.UUID, seasonNumber int32, airDate pgtype.Date) (dbgen.MediaSeason, error) {
	return r.Q.UpsertSeason(ctx, dbgen.UpsertSeasonParams{
		MediaItemID:  mediaItemID,
		SeasonNumber: seasonNumber,
		AirDate:      airDate,
	})
}

func (r *Repository) ListEpisodesForSeason(ctx context.Context, seasonID pgtype.UUID) ([]dbgen.MediaEpisode, error) {
	return r.Q.ListEpisodesForSeason(ctx, seasonID)
}

func (r *Repository) UpsertEpisode(ctx context.Context, seasonID pgtype.UUID, episodeNumber int32, title *string, airDate pgtype.Date, tmdbID *int64, tvdbID *int64) (dbgen.MediaEpisode, error) {
	return r.Q.UpsertEpisode(ctx, dbgen.UpsertEpisodeParams{
		SeasonID:      seasonID,
		EpisodeNumber: episodeNumber,
		Title:         title,
		AirDate:       airDate,
		TmdbID:        tmdbID,
		TvdbID:        tvdbID,
	})
}

func (r *Repository) GetMediaFileByPath(ctx context.Context, path string) (dbgen.MediaFile, error) {
	return r.Q.GetMediaFileByPath(ctx, path)
}

func (r *Repository) CreateMediaFile(ctx context.Context, mediaItemID pgtype.UUID, seasonID, episodeID *pgtype.UUID, path string) (dbgen.MediaFile, error) {
	var season, episode pgtype.UUID
	if seasonID != nil {
		season = *seasonID
	} // else zero value => NULL
	if episodeID != nil {
		episode = *episodeID
	}
	return r.Q.CreateMediaFile(ctx, dbgen.CreateMediaFileParams{
		MediaItemID: mediaItemID,
		SeasonID:    season,
		EpisodeID:   episode,
		Path:        path,
	})
}

func (r *Repository) DeleteMediaFile(ctx context.Context, id pgtype.UUID) error {
	return r.Q.DeleteMediaFile(ctx, id)
}
