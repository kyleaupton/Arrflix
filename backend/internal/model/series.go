package model

type Series struct {
	TmdbID      int64  `json:"tmdbId"`
	Title       string `json:"title"`
	PosterPath  string `json:"posterPath"`
	ReleaseDate string `json:"releaseDate"`
}
