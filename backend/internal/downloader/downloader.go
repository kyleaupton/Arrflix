package downloader

import (
	"context"
	"errors"
	"time"
)

var ErrUnsupported = errors.New("operation unsupported")

type Type string

const (
	TypeQbittorrent Type = "qbittorrent"
	TypeSABnzbd     Type = "sabnzbd" // future
	TypeNZBGet      Type = "nzbget"  // future
)

type InstanceID string // your DB UUID string, etc.

type AddRequest struct {
	// One of these is required depending on kind
	MagnetURL string // torrent
	NZBURL    string // usenet (optional future)
	// Maybe also: TorrentFileBytes []byte

	Category string
	Tags     []string

	// Optional behavior knobs (some clients ignore)
	Paused bool

	// If set, downloader should try to download here (if supported)
	SavePath string
}

type AddResult struct {
	ExternalID string // torrent hash / nzb id / whatever
	Name       string // best-effort
}

type JobStatus string

const (
	StatusUnknown     JobStatus = "unknown"
	StatusQueued      JobStatus = "queued"
	StatusDownloading JobStatus = "downloading"
	StatusCompleted   JobStatus = "completed"
	StatusSeeding     JobStatus = "seeding"
	StatusPaused      JobStatus = "paused"
	StatusErrored     JobStatus = "errored"
)

type Item struct {
	ExternalID string

	Name     string
	Status   JobStatus
	Progress float64 // 0..1

	// Best-effort locations as seen by the client
	SavePath    string
	ContentPath string

	AddedAt time.Time
}

type File struct {
	Path     string // relative to root/content path if possible
	Size     int64
	Progress float64 // 0..1 if available
	Priority int     // optional
}

type Client interface {
	Type() Type
	InstanceID() InstanceID

	// Add a download (magnet/NZB/etc)
	Add(ctx context.Context, req AddRequest) (AddResult, error)

	// Fetch a single item by external id (infohash, etc)
	Get(ctx context.Context, externalID string) (Item, error)

	// List all items (optional - can return ErrUnsupported)
	List(ctx context.Context) ([]Item, error)

	// Optional but very useful for import selection
	ListFiles(ctx context.Context, externalID string) ([]File, error)

	// Housekeeping knobs (optional to implement; you can return ErrUnsupported)
	Pause(ctx context.Context, externalID string) error
	Resume(ctx context.Context, externalID string) error
	Remove(ctx context.Context, externalID string, deleteData bool) error
}

type ConfigRecord struct {
	ID       InstanceID
	Type     Type
	URL      string
	Username *string
	Password *string
	Config   []byte // JSON from DB (type-specific config)
}

type Builder func(rec ConfigRecord) (Client, error)
