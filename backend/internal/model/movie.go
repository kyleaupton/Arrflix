package model

type Movie struct {
	// General movie details
	TmdbID      int64  `json:"tmdbId"`
	Title       string `json:"title"`
	ReleaseDate string `json:"releaseDate"`
}
