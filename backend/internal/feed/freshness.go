package feed

import (
	"sync"
	"time"

	"github.com/kyleaupton/snaggle/backend/internal/model"
)

// FreshnessTracker tracks when rows were last shown to enforce rotation
type FreshnessTracker interface {
	GetFreshnessFactor(intent model.RowIntent) float64
	RecordShown(intents []model.RowIntent)
}

// InMemoryFreshnessTracker implements FreshnessTracker using in-memory storage
type InMemoryFreshnessTracker struct {
	mu        sync.RWMutex
	lastShown map[model.RowIntent]time.Time
}

// NewInMemoryFreshnessTracker creates a new freshness tracker
func NewInMemoryFreshnessTracker() *InMemoryFreshnessTracker {
	return &InMemoryFreshnessTracker{
		lastShown: make(map[model.RowIntent]time.Time),
	}
}

// GetFreshnessFactor returns a decay multiplier based on when the row was last shown
// - Rows shown in last request: freshness_factor = 0.3
// - Rows shown 1-6 hours ago: freshness_factor = 0.7
// - Rows shown 6-24 hours ago: freshness_factor = 0.9
// - Rows not shown in 24+ hours: freshness_factor = 1.0
func (t *InMemoryFreshnessTracker) GetFreshnessFactor(intent model.RowIntent) float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	lastTime, exists := t.lastShown[intent]
	if !exists {
		return 1.0 // never shown before
	}

	elapsed := time.Since(lastTime)

	switch {
	case elapsed < 1*time.Hour:
		return 0.3 // just shown
	case elapsed < 6*time.Hour:
		return 0.7 // shown recently
	case elapsed < 24*time.Hour:
		return 0.9 // shown today
	default:
		return 1.0 // fresh again
	}
}

// RecordShown updates the last shown timestamp for the given row intents
func (t *InMemoryFreshnessTracker) RecordShown(intents []model.RowIntent) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	for _, intent := range intents {
		t.lastShown[intent] = now
	}
}
