package model

import "github.com/google/uuid"

type Plan struct {
	// this is given to us
	Torrent Torrent // what to download

	// these are determined by the policy engine
	DownloaderID   string // how to download
	LibraryID      string // where to move/hardlink/copy the file to
	NameTemplateID string // how to name the file
}

type Torrent struct {
	ID        string
	IndexerID string
}

// TorrentMetadata represents metadata about a torrent for policy evaluation
type TorrentMetadata struct {
	Size       uint64   // file size in bytes
	Seeders    uint     // number of seeders
	Peers      uint     // number of peers
	Title      string   // torrent title
	Tracker    string   // tracker/indexer name
	TrackerID  string   // tracker/indexer ID
	Categories []string // category tags
}

type Policy struct {
	ID          uuid.UUID
	Name        string
	Description string
	Enabled     bool
	Priority    int

	Condition Rule
	Actions   []Action
}

type Operator string

const (
	OpEq       Operator = "=="
	OpNe       Operator = "!="
	OpGt       Operator = ">"
	OpGte      Operator = ">="
	OpLt       Operator = "<"
	OpLte      Operator = "<="
	OpContains Operator = "contains"
	OpIn       Operator = "in"
	OpNotIn    Operator = "not in"
	OpAnd      Operator = "and"
	OpOr       Operator = "or"
	OpNot      Operator = "not"
)

type Rule struct {
	ID       uuid.UUID
	Left     string
	Operator Operator
	Right    string
}

type ActionType string

const (
	ActionSetDownloader   ActionType = "set_downloader"
	ActionSetLibrary      ActionType = "set_library"
	ActionSetNameTemplate ActionType = "set_name_template"
	ActionStopProcessing  ActionType = "stop_processing"
)

type Action struct {
	ID    uuid.UUID
	Type  ActionType
	Value string
}
