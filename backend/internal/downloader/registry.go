package downloader

import "fmt"

// Registry manages builder functions for different downloader types
type Registry struct {
	builders map[Type]Builder
}

// NewRegistry creates a new registry
func NewRegistry() *Registry {
	return &Registry{builders: map[Type]Builder{}}
}

// Register registers a builder function for a downloader type
func (r *Registry) Register(t Type, b Builder) {
	r.builders[t] = b
}

// Build builds a client instance from a config record
func (r *Registry) Build(rec ConfigRecord) (Client, error) {
	b, ok := r.builders[rec.Type]
	if !ok {
		return nil, fmt.Errorf("unknown downloader type: %s", rec.Type)
	}
	return b(rec)
}
