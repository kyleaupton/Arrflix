package indexer

import "context"

// IndexerSource abstracts indexer search operations.
// Implementations handle protocol-specific quirks and validate results at the boundary.
type IndexerSource interface {
	// Search performs a search query and returns validated results.
	// Invalid results (e.g., missing download URL) are filtered out.
	Search(ctx context.Context, query SearchQuery) ([]SearchResult, error)

	// ListIndexers returns information about all configured indexers.
	ListIndexers(ctx context.Context) ([]IndexerInfo, error)

	// Test verifies connectivity to the indexer backend.
	Test(ctx context.Context) error
}
