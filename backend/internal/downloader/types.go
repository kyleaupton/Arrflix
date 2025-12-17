package downloader

// Protocol represents the download protocol (torrent or usenet)
type Protocol string

const (
	ProtocolTorrent Protocol = "torrent"
	ProtocolUsenet  Protocol = "usenet"
)

// Config represents type-specific configuration stored in config_json
// This is a base struct that can be embedded by specific downloader types
type Config struct {
	// Common fields can be added here if needed
	// Type-specific fields should be defined in their own config structs
}
