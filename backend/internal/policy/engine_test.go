package policy

import (
	"context"
	"testing"

	"github.com/kyleaupton/snaggle/backend/internal"
	"github.com/kyleaupton/snaggle/backend/internal/model"
)

func TestEngine_Evaluate(t *testing.T) {
	r := internal.GetRepo()
	engine := NewEngine(r)

	plan, err := engine.Evaluate(context.Background(), model.EvaluateParams{
		TorrentURL: "https://example.com/torrent.torrent",
		Metadata: model.TorrentMetadata{
			Size:       1000,
			Seeders:    10,
			Peers:      10,
			Title:      "Test Torrent",
			Tracker:    "Test Tracker",
			TrackerID:  "1234567890",
			Categories: []string{"Test Category"},
		},
		MediaType: model.MediaTypeMovie,
	})

	if err != nil {
		t.Fatalf("error evaluating plan: %v", err)
	}

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
