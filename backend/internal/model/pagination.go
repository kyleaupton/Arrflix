package model

// Pagination contains metadata for paginated responses
type Pagination struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
}

// PaginatedLibraryResponse is the envelope for paginated library items
// Note: Using concrete type instead of generics for Swagger compatibility
type PaginatedLibraryResponse struct {
	Data       []LibraryItem `json:"data"`
	Pagination Pagination    `json:"pagination"`
}

// LibraryItem is the enriched media item returned by the library endpoint
type LibraryItem struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	Year       *int32 `json:"year,omitempty"`
	TmdbID     *int64 `json:"tmdbId,omitempty"`
	PosterPath string `json:"posterPath,omitempty"`
	CreatedAt  string `json:"createdAt"`
}

