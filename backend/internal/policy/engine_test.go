package policy

import (
	"context"
	"os"
	"testing"

	"github.com/kyleaupton/snaggle/backend/internal"
	"github.com/kyleaupton/snaggle/backend/internal/logger"
	"github.com/kyleaupton/snaggle/backend/internal/model"
	"github.com/kyleaupton/snaggle/backend/internal/release"
)

func TestEngine_Evaluate(t *testing.T) {
	if os.Getenv("SNAGGLE_INTEGRATION") != "1" {
		t.Skip("skipping integration test (set SNAGGLE_INTEGRATION=1 to enable)")
	}

	r := internal.GetRepo()
	logg := logger.New(true)
	engine := NewEngine(r, logg)

	candidate := model.DownloadCandidate{
		Protocol:  "http",
		Filename:  "test.torrent",
		Link:      "https://example.com/torrent.torrent",
		Indexer:   "Test Indexer",
		IndexerID: 1234567890,
		GUID:      "1234567890",
		Peers:     10,
		Seeders:   10,
		Age:       1000,
		AgeHours:  10,
	}
	q := release.Parse(candidate.Title)
	evalCtx := model.NewEvaluationContext(candidate, q)

	trace, err := engine.Evaluate(context.Background(), evalCtx)

	if err != nil {
		t.Fatalf("error evaluating plan: %v", err)
	}

	plan := trace.FinalPlan
	if plan.DownloaderID != "Test Downloader" {
		t.Fatalf("expected downloader ID to be Test Downloader, got %s", plan.DownloaderID)
	}

	if plan.LibraryID != "Test Library" {
		t.Fatalf("expected library ID to be Test Library, got %s", plan.LibraryID)
	}

	if plan.NameTemplateID != "Test Name Template" {
		t.Fatalf("expected name template ID to be Test Name Template, got %s", plan.NameTemplateID)
	}
}
