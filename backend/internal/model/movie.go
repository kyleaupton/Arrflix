package model

type Movie struct {
	// General movie details
	TmdbID       int64  `json:"tmdbId"`
	Title        string `json:"title"`
	PosterPath   string `json:"posterPath"`
	BackdropPath string `json:"backdropPath"`
	ReleaseDate  string `json:"releaseDate"`
}
