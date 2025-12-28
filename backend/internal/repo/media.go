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
	GetMediaItemByTmdbIDAndType(ctx context.Context, tmdbID int64, typ string) (dbgen.MediaItem, error)
	CreateMediaItem(ctx context.Context, typ, title string, year *int32, tmdbID *int32) (dbgen.MediaItem, error)
	UpdateMediaItem(ctx context.Context, id pgtype.UUID, title string, year *int32, tmdbID *int32) (dbgen.MediaItem, error)
	DeleteMediaItem(ctx context.Context, id pgtype.UUID) error

	// Seasons
	ListSeasonsForMedia(ctx context.Context, mediaItemID pgtype.UUID) ([]dbgen.MediaSeason, error)
	GetSeason(ctx context.Context, id pgtype.UUID) (dbgen.MediaSeason, error)
	GetSeasonByNumber(ctx context.Context, mediaItemID pgtype.UUID, seasonNumber int32) (dbgen.MediaSeason, error)
	UpsertSeason(ctx context.Context, mediaItemID pgtype.UUID, seasonNumber int32, airDate pgtype.Date) (dbgen.MediaSeason, error)

	// Episodes
	ListEpisodesForSeason(ctx context.Context, seasonID pgtype.UUID) ([]dbgen.MediaEpisode, error)
	GetEpisode(ctx context.Context, id pgtype.UUID) (dbgen.MediaEpisode, error)
	GetEpisodeByNumber(ctx context.Context, seasonID pgtype.UUID, episodeNumber int32) (dbgen.MediaEpisode, error)
	UpsertEpisode(ctx context.Context, seasonID pgtype.UUID, episodeNumber int32, title *string, airDate pgtype.Date, tmdbID *int64, tvdbID *int64) (dbgen.MediaEpisode, error)

	// Files
	GetMediaFileByLibraryAndPath(ctx context.Context, libraryID pgtype.UUID, path string) (dbgen.MediaFile, error)
	CreateMediaFile(ctx context.Context, libraryID, mediaItemID pgtype.UUID, seasonID, episodeID *pgtype.UUID, path string, status *string) (dbgen.MediaFile, error)
	ListMediaFilesForItem(ctx context.Context, mediaItemID pgtype.UUID) ([]dbgen.ListMediaFilesForItemRow, error)
	ListEpisodeAvailabilityForSeries(ctx context.Context, mediaItemID pgtype.UUID) ([]dbgen.ListEpisodeAvailabilityForSeriesRow, error)
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

func (r *Repository) GetMediaItemByTmdbIDAndType(ctx context.Context, tmdbID int64, typ string) (dbgen.MediaItem, error) {
	return r.Q.GetMediaItemByTmdbIDAndType(ctx, dbgen.GetMediaItemByTmdbIDAndTypeParams{
		TmdbID: &tmdbID,
		Type:   typ,
	})
}

func (r *Repository) CreateMediaItem(ctx context.Context, typ, title string, year *int32, tmdbID *int64) (dbgen.MediaItem, error) {
	return r.Q.CreateMediaItem(ctx, dbgen.CreateMediaItemParams{
		Type:   typ,
		Title:  title,
		Year:   year,
		TmdbID: tmdbID,
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

func (r *Repository) GetSeason(ctx context.Context, id pgtype.UUID) (dbgen.MediaSeason, error) {
	return r.Q.GetSeason(ctx, id)
}

func (r *Repository) GetSeasonByNumber(ctx context.Context, mediaItemID pgtype.UUID, seasonNumber int32) (dbgen.MediaSeason, error) {
	return r.Q.GetSeasonByNumber(ctx, dbgen.GetSeasonByNumberParams{
		MediaItemID:  mediaItemID,
		SeasonNumber: seasonNumber,
	})
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

func (r *Repository) GetEpisode(ctx context.Context, id pgtype.UUID) (dbgen.MediaEpisode, error) {
	return r.Q.GetEpisode(ctx, id)
}

func (r *Repository) GetEpisodeByNumber(ctx context.Context, seasonID pgtype.UUID, episodeNumber int32) (dbgen.MediaEpisode, error) {
	return r.Q.GetEpisodeByNumber(ctx, dbgen.GetEpisodeByNumberParams{
		SeasonID:      seasonID,
		EpisodeNumber: episodeNumber,
	})
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

func (r *Repository) GetMediaFileByLibraryAndPath(ctx context.Context, libraryID pgtype.UUID, path string) (dbgen.MediaFile, error) {
	return r.Q.GetMediaFileByLibraryAndPath(ctx, dbgen.GetMediaFileByLibraryAndPathParams{
		LibraryID: libraryID,
		Path:      path,
	})
}

func (r *Repository) CreateMediaFile(ctx context.Context, libraryID, mediaItemID pgtype.UUID, seasonID, episodeID *pgtype.UUID, path string, status *string) (dbgen.MediaFile, error) {
	var season, episode pgtype.UUID
	if seasonID != nil {
		season = *seasonID
	} // else zero value => NULL
	if episodeID != nil {
		episode = *episodeID
	}
	var st *string
	if status != nil {
		st = status
	}
	return r.Q.CreateMediaFile(ctx, dbgen.CreateMediaFileParams{
		LibraryID:   libraryID,
		MediaItemID: mediaItemID,
		SeasonID:    season,
		EpisodeID:   episode,
		Path:        path,
		Status:      st,
	})
}

func (r *Repository) ListMediaFilesForItem(ctx context.Context, mediaItemID pgtype.UUID) ([]dbgen.ListMediaFilesForItemRow, error) {
	return r.Q.ListMediaFilesForItem(ctx, mediaItemID)
}

func (r *Repository) ListEpisodeAvailabilityForSeries(ctx context.Context, mediaItemID pgtype.UUID) ([]dbgen.ListEpisodeAvailabilityForSeriesRow, error) {
	return r.Q.ListEpisodeAvailabilityForSeries(ctx, mediaItemID)
}

func (r *Repository) DeleteMediaFile(ctx context.Context, id pgtype.UUID) error {
	return r.Q.DeleteMediaFile(ctx, id)
}
