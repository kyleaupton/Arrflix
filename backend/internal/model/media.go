package model

type MovieRail struct {
	TmdbID      int64   `json:"tmdbId"`
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	PosterPath  string  `json:"posterPath"`
	ReleaseDate string  `json:"releaseDate"`
	Year        *int32  `json:"year,omitempty"`
	Genres      []int64 `json:"genres,omitempty"`
	Tagline     string  `json:"tagline,omitempty"`
	IsInLibrary bool    `json:"isInLibrary"`
}

type Movie struct {
	TmdbID int64 `json:"tmdbId"`

	Title    string `json:"title"`
	Overview string `json:"overview"`
	Tagline  string `json:"tagline"`

	// Stats
	Status              string              `json:"status"`
	ReleaseDate         string              `json:"releaseDate"`
	Runtime             int                 `json:"runtime"`
	OriginalLanguage    string              `json:"originalLanguage"`
	OriginCountry       []string            `json:"originCountry"`
	ProductionCompanies []ProductionCompany `json:"productionCompanies"`
	ProductionCountries []ProductionCountry `json:"productionCountries"`

	PosterPath   string `json:"posterPath"`
	BackdropPath string `json:"backdropPath"`
}

type SeriesRail struct {
	TmdbID      int64   `json:"tmdbId"`
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	PosterPath  string  `json:"posterPath"`
	ReleaseDate string  `json:"releaseDate"`
	Year        *int32  `json:"year,omitempty"`
	Genres      []int64 `json:"genres,omitempty"`
	Tagline     string  `json:"tagline,omitempty"`
	IsInLibrary bool    `json:"isInLibrary"`
}

type Series struct {
	TmdbID int64 `json:"tmdbId"`

	Title    string `json:"title"`
	Overview string `json:"overview"`
	Tagline  string `json:"tagline"`
	Status   string `json:"status"`

	Seasons []Season `json:"seasons"`

	PosterPath   string `json:"posterPath"`
	BackdropPath string `json:"backdropPath"`

	FirstAirDate string `json:"firstAirDate"`
	LastAirDate  string `json:"lastAirDate"`
	InProduction bool   `json:"inProduction"`

	CreatedBy []Person  `json:"createdBy"`
	Networks  []Network `json:"networks"`
	Genres    []Genre   `json:"genres"`
}

type Season struct {
	TmdbID       int64  `json:"tmdbId"`
	SeasonNumber int64  `json:"seasonNumber"`
	Title        string `json:"title"`
	Overview     string `json:"overview"`
	PosterPath   string `json:"posterPath"`
	AirDate      string `json:"airDate"`
}

type Person struct {
	TmdbID      int64  `json:"tmdbId"`
	Name        string `json:"name"`
	ProfilePath string `json:"profilePath"`
}

type Network struct {
	TmdbID   int64  `json:"tmdbId"`
	Name     string `json:"name"`
	LogoPath string `json:"logoPath"`
}

type Genre struct {
	TmdbID int64  `json:"tmdbId"`
	Name   string `json:"name"`
}

type ProductionCompany struct {
	TmdbID        int64  `json:"tmdbId"`
	Name          string `json:"name"`
	LogoPath      string `json:"logoPath"`
	OriginCountry string `json:"originCountry"`
}

type ProductionCountry struct {
	Iso3166_1 string `json:"iso3166_1"`
	Name      string `json:"name"`
}

type MediaType string

const (
	MediaTypeMovie  MediaType = "movie"
	MediaTypeSeries MediaType = "series"
)
