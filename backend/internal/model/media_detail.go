package model

type FileInfo struct {
	ID            string `json:"id"`
	LibraryID     string `json:"libraryId"`
	Path          string `json:"path"` // relative to library root
	Status        string `json:"status"`
	SeasonNumber  *int32 `json:"seasonNumber,omitempty"`
	EpisodeNumber *int32 `json:"episodeNumber,omitempty"`
}

type LibraryAvailability struct {
	LibraryID    string `json:"libraryId"`
	FileCount    int    `json:"fileCount"`
	StatusRollup string `json:"statusRollup"`
}

type Availability struct {
	IsInLibrary bool                  `json:"isInLibrary"`
	Libraries   []LibraryAvailability `json:"libraries"`
}

type MovieDetail struct {
	TmdbID int64  `json:"tmdbId"`
	Title  string `json:"title"`
	Year   *int32 `json:"year,omitempty"`

	Overview     string  `json:"overview"`
	Tagline      string  `json:"tagline,omitempty"`
	Status       string  `json:"status"`
	ReleaseDate  string  `json:"releaseDate,omitempty"`
	Runtime      int     `json:"runtime,omitempty"`
	Genres       []Genre `json:"genres,omitempty"`
	PosterPath   string  `json:"posterPath,omitempty"`
	BackdropPath string  `json:"backdropPath,omitempty"`

	Files   []FileInfo `json:"files"`
	Credits *Credits   `json:"credits,omitempty"`
	Videos  []Video    `json:"videos,omitempty"`
}

type EpisodeAvailability struct {
	SeasonNumber  int32   `json:"seasonNumber"`
	EpisodeNumber int32   `json:"episodeNumber"`
	Title         *string `json:"title,omitempty"`
	AirDate       *string `json:"airDate,omitempty"`
	Available     bool    `json:"available"`
	FileID        *string `json:"fileId,omitempty"`
}

type SeasonDetail struct {
	SeasonNumber int32                 `json:"seasonNumber"`
	Episodes     []EpisodeAvailability `json:"episodes"`
}

type SeriesDetail struct {
	TmdbID int64  `json:"tmdbId"`
	Title  string `json:"title"`
	Year   *int32 `json:"year,omitempty"`

	Overview     string  `json:"overview"`
	Tagline      string  `json:"tagline,omitempty"`
	Status       string  `json:"status"`
	FirstAirDate string  `json:"firstAirDate,omitempty"`
	LastAirDate  string  `json:"lastAirDate,omitempty"`
	InProduction bool    `json:"inProduction"`
	Genres       []Genre `json:"genres,omitempty"`
	PosterPath   string  `json:"posterPath,omitempty"`
	BackdropPath string  `json:"backdropPath,omitempty"`

	Availability Availability   `json:"availability"`
	Files        []FileInfo     `json:"files"`
	Seasons      []SeasonDetail `json:"seasons"`
	Credits      *Credits       `json:"credits,omitempty"`
	Videos       []Video        `json:"videos,omitempty"`
}
