package model

type Rail struct {
	ID     string       `json:"id"`
	Title  string       `json:"title"`
	Type   string       `json:"type"`
	Series []SeriesRail `json:"series"`
	Movies []MovieRail  `json:"movies"`
}
