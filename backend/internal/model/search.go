package model

type SearchResult struct {
	ID          int64   `json:"id"`
	MediaType   string  `json:"mediaType"` // "movie", "tv", "person"
	Title       string  `json:"title"`
	PosterPath  *string `json:"posterPath,omitempty"`
	Year        *int    `json:"year,omitempty"`
	Overview    *string `json:"overview,omitempty"`
	IsInLibrary bool    `json:"isInLibrary"`
}

type SearchResponse struct {
	Results      []SearchResult `json:"results"`
	TotalResults int            `json:"totalResults"`
	Query        string         `json:"query"`
}
