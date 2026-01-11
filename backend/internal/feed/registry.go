package feed

import (
	"context"
	"sort"

	"github.com/kyleaupton/Arrflix/internal/model"
)

// Registry holds all available row definitions
type Registry struct {
	rows map[model.RowIntent]model.RowDefinition
}

// NewRegistry creates a registry with built-in row definitions
func NewRegistry() *Registry {
	r := &Registry{rows: make(map[model.RowIntent]model.RowDefinition)}
	r.registerBuiltinRows()
	return r
}

func (r *Registry) registerBuiltinRows() {
	defs := []model.RowDefinition{
		// 1. Trending This Week
		{
			Intent:      model.IntentTrendingWeek,
			Title:       "Trending This Week",
			ContentKind: model.ContentKindMixed,
			Sources: []model.SourceConfig{
				{Provider: "tmdb", Endpoint: "trending", Params: map[string]string{"media_type": "movie", "time_window": "day"}},
				{Provider: "tmdb", Endpoint: "trending", Params: map[string]string{"media_type": "tv", "time_window": "day"}},
			},
			TargetSize: 20,
			FetchSize:  40,
			Ranking:    model.RankingDefault,
			Diversity:  &model.DiversityRules{MaxPerGenre: 3},
			Weight:     1.0,
		},
		// 2. New & Noteworthy
		{
			Intent:      model.IntentNewAndNoteworthy,
			Title:       "New & Noteworthy",
			Subtitle:    "Recent releases getting attention",
			ContentKind: model.ContentKindMixed,
			Sources: []model.SourceConfig{
				{Provider: "tmdb", Endpoint: "now_playing"},
				{Provider: "tmdb", Endpoint: "on_the_air"},
			},
			TargetSize: 20,
			FetchSize:  40,
			Ranking:    model.RankingPopularity,
			Diversity:  &model.DiversityRules{MaxPerGenre: 2, MaxPerLanguage: 5},
			Weight:     0.9,
		},
		// 3. Hidden Gems
		{
			Intent:      model.IntentHiddenGems,
			Title:       "Hidden Gems",
			Subtitle:    "Overlooked favorites worth discovering",
			ContentKind: model.ContentKindMixed,
			Sources: []model.SourceConfig{
				{Provider: "tmdb", Endpoint: "discover", Params: map[string]string{
					"media_type":       "movie",
					"vote_average.gte": "7.5",
					"vote_count.lte":   "500",
				}},
				{Provider: "tmdb", Endpoint: "discover", Params: map[string]string{
					"media_type":       "tv",
					"vote_average.gte": "7.5",
					"vote_count.lte":   "500",
				}},
			},
			TargetSize: 20,
			FetchSize:  40,
			Ranking:    model.RankingRating,
			Diversity:  &model.DiversityRules{MaxPerGenre: 2},
			Weight:     0.5,
		},
		// 4. Critically Acclaimed
		{
			Intent:      model.IntentCriticallyAcclaimed,
			Title:       "Critically Acclaimed",
			ContentKind: model.ContentKindMixed,
			Sources: []model.SourceConfig{
				{Provider: "tmdb", Endpoint: "top_rated", Params: map[string]string{"media_type": "movie"}},
				{Provider: "tmdb", Endpoint: "top_rated", Params: map[string]string{"media_type": "tv"}},
			},
			TargetSize: 20,
			FetchSize:  40,
			Ranking:    model.RankingRating,
			Diversity:  &model.DiversityRules{MaxPerGenre: 3},
			Weight:     0.6,
		},
		// 5. Just Released
		{
			Intent:      model.IntentJustReleased,
			Title:       "Just Released",
			Subtitle:    "Fresh content from this month",
			ContentKind: model.ContentKindMixed,
			Sources: []model.SourceConfig{
				{Provider: "tmdb", Endpoint: "discover", Params: map[string]string{
					"media_type":           "movie",
					"primary_release_date": "recent30",
				}},
				{Provider: "tmdb", Endpoint: "discover", Params: map[string]string{
					"media_type":     "tv",
					"first_air_date": "recent30",
				}},
			},
			TargetSize: 20,
			FetchSize:  40,
			Ranking:    model.RankingRecent,
			Diversity:  &model.DiversityRules{MaxPerGenre: 3},
			Weight:     0.8,
		},
		// 6. Coming Soon
		{
			Intent:      model.IntentComingSoon,
			Title:       "Coming Soon",
			Subtitle:    "Get ready for these upcoming releases",
			ContentKind: model.ContentKindMovie,
			Sources: []model.SourceConfig{
				{Provider: "tmdb", Endpoint: "upcoming"},
			},
			TargetSize: 20,
			FetchSize:  40,
			Ranking:    model.RankingRecent,
			Diversity:  &model.DiversityRules{MaxPerGenre: 2},
			Weight:     0.4,
		},
		// 7. Binge-Worthy Series
		{
			Intent:      model.IntentBingeWorthy,
			Title:       "Binge-Worthy Series",
			Subtitle:    "Series you can't stop watching",
			ContentKind: model.ContentKindTV,
			Sources: []model.SourceConfig{
				{Provider: "tmdb", Endpoint: "popular", Params: map[string]string{"media_type": "tv"}},
			},
			TargetSize: 20,
			FetchSize:  40,
			Ranking:    model.RankingPopularity,
			Diversity:  &model.DiversityRules{MaxPerGenre: 2},
			Weight:     0.7,
		},
		// 8. Fan Favorites
		{
			Intent:      model.IntentFanFavorites,
			Title:       "Fan Favorites",
			Subtitle:    "Beloved by audiences everywhere",
			ContentKind: model.ContentKindMixed,
			Sources: []model.SourceConfig{
				{Provider: "tmdb", Endpoint: "discover", Params: map[string]string{
					"media_type":       "movie",
					"vote_count.gte":   "1000",
					"vote_average.gte": "8.0",
				}},
				{Provider: "tmdb", Endpoint: "discover", Params: map[string]string{
					"media_type":       "tv",
					"vote_count.gte":   "1000",
					"vote_average.gte": "8.0",
				}},
			},
			TargetSize: 20,
			FetchSize:  40,
			Ranking:    model.RankingPopularity,
			Diversity:  &model.DiversityRules{MaxPerGenre: 3},
			Weight:     0.5,
		},
	}

	for _, d := range defs {
		r.rows[d.Intent] = d
	}
}

// SelectRows returns eligible row definitions based on context
func (r *Registry) SelectRows(ctx context.Context, hasSignals bool) []model.RowDefinition {
	var selected []model.RowDefinition

	for _, def := range r.rows {
		// Skip rows that require user signals when signals aren't available
		if def.RequiresSignal && !hasSignals {
			continue
		}
		selected = append(selected, def)
	}

	return selected
}

// RowScore represents a row with its calculated score
type RowScore struct {
	Definition model.RowDefinition
	Score      float64
}

// ScoreAndOrder applies dynamic scoring and returns ordered rows
func (r *Registry) ScoreAndOrder(rows []model.RowDefinition, freshness FreshnessTracker, hasSignals bool) []model.RowDefinition {
	scores := make([]RowScore, 0, len(rows))

	for _, row := range rows {
		// Calculate freshness factor
		freshnessFactor := freshness.GetFreshnessFactor(row.Intent)

		// Calculate signal boost (1.0 for now, can be enhanced based on user signals)
		signalBoost := 1.0
		if hasSignals && row.RequiresSignal {
			signalBoost = 1.5
		}

		// Final score = weight × freshness × signal boost
		score := row.Weight * freshnessFactor * signalBoost

		scores = append(scores, RowScore{
			Definition: row,
			Score:      score,
		})
	}

	// Sort by score descending (higher scores first)
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	// Extract ordered definitions
	ordered := make([]model.RowDefinition, len(scores))
	for i, rs := range scores {
		ordered[i] = rs.Definition
	}

	return ordered
}
