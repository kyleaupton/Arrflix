package model

type Series struct {
	TmdbID      int64  `json:"tmdbId"`
	Title       string `json:"title"`
	ReleaseDate string `json:"releaseDate"`
}
