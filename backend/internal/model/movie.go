package model

type Movie struct {
	// General movie details
	TmdbID      int64  `json:"tmdbId"`
	Title       string `json:"title"`
	PosterPath  string `json:"posterPath"`
	ReleaseDate string `json:"releaseDate"`
}
