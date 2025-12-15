package model

import "time"

// DownloadCandidate represents a download candidate from Prowlarr search results
type DownloadCandidate struct {
	Protocol    string    `json:"protocol"`
	Filename    string    `json:"filename"`
	Link        string    `json:"link"`
	Indexer     string    `json:"indexer"`
	IndexerID   int64     `json:"indexerId"`
	GUID        string    `json:"guid"`
	Peers       int       `json:"peers"`   // leechers
	Seeders     int       `json:"seeders"` // seeders
	Age         int64     `json:"age"`     // age in seconds
	AgeHours    float64   `json:"ageHours"`
	Size        int64     `json:"size"` // size in bytes
	Grabs       int       `json:"grabs"`
	Categories  []string  `json:"categories"`
	PublishDate time.Time `json:"publishDate"`
	Title       string    `json:"title"`
}

// EnqueueCandidateRequest is the request body for enqueueing a download candidate
type EnqueueCandidateRequest struct {
	IndexerID int64  `json:"indexerId"`
	GUID      string `json:"guid"`
}
