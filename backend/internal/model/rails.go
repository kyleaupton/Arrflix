package model

type Rail struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Type   string   `json:"type"`
	Series []Series `json:"series"`
	Movies []Movie  `json:"movies"`
}
