package jackett

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetIndexerConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		indexerID      string
		responseStatus int
		responseBody   string
		wantError      bool
		wantConfig     *IndexerConfig
	}{
		{
			name:           "successful config retrieval",
			indexerID:      "test-indexer",
			responseStatus: http.StatusOK,
			responseBody: `{
				"id": "test-indexer",
				"name": "Test Indexer",
				"description": "A test indexer",
				"type": "public",
				"configured": true,
				"enabled": true,
				"fields": {
					"username": "testuser",
					"password": "testpass"
				}
			}`,
			wantError: false,
			wantConfig: &IndexerConfig{
				ID:          "test-indexer",
				Name:        "Test Indexer",
				Description: "A test indexer",
				Type:        "public",
				Configured:  true,
				Enabled:     true,
				Fields: map[string]interface{}{
					"username": "testuser",
					"password": "testpass",
				},
			},
		},
		{
			name:           "indexer not found",
			indexerID:      "nonexistent",
			responseStatus: http.StatusNotFound,
			responseBody: `{
				"error": "Indexer not found"
			}`,
			wantError: true,
		},
		{
			name:           "server error",
			indexerID:      "test-indexer",
			responseStatus: http.StatusInternalServerError,
			responseBody: `{
				"error": "Internal server error"
			}`,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("expected GET request, got %s", r.Method)
				}
				expectedPath := "/api/v2.0/indexers/" + tt.indexerID + "/config"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}
				if r.URL.Query().Get("apikey") != testAPIKey {
					t.Errorf("expected API key %s, got %s", testAPIKey, r.URL.Query().Get("apikey"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.responseStatus)
				w.Write([]byte(tt.responseBody))
			}))
			defer srv.Close()

			client, err := New(Settings{
				ApiURL: srv.URL,
				ApiKey: testAPIKey,
				Client: srv.Client(),
			})
			if err != nil {
				t.Fatal(err)
			}

			config, err := client.GetIndexerConfig(context.Background(), tt.indexerID)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if config == nil {
				t.Errorf("expected config, got nil")
				return
			}

			if config.ID != tt.wantConfig.ID {
				t.Errorf("expected ID %s, got %s", tt.wantConfig.ID, config.ID)
			}
			if config.Name != tt.wantConfig.Name {
				t.Errorf("expected Name %s, got %s", tt.wantConfig.Name, config.Name)
			}
			if config.Description != tt.wantConfig.Description {
				t.Errorf("expected Description %s, got %s", tt.wantConfig.Description, config.Description)
			}
			if config.Type != tt.wantConfig.Type {
				t.Errorf("expected Type %s, got %s", tt.wantConfig.Type, config.Type)
			}
			if config.Configured != tt.wantConfig.Configured {
				t.Errorf("expected Configured %v, got %v", tt.wantConfig.Configured, config.Configured)
			}
			if config.Enabled != tt.wantConfig.Enabled {
				t.Errorf("expected Enabled %v, got %v", tt.wantConfig.Enabled, config.Enabled)
			}
		})
	}
}

func TestSaveIndexerConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		indexerID      string
		config         *IndexerConfigRequest
		responseStatus int
		responseBody   string
		wantError      bool
		wantConfig     *IndexerConfig
	}{
		{
			name:      "successful config save",
			indexerID: "test-indexer",
			config: &IndexerConfigRequest{
				Name:        "Updated Test Indexer",
				Description: "An updated test indexer",
				Type:        "private",
				Enabled:     true,
				Fields: map[string]interface{}{
					"username": "newuser",
					"password": "newpass",
				},
			},
			responseStatus: http.StatusOK,
			responseBody: `{
				"id": "test-indexer",
				"name": "Updated Test Indexer",
				"description": "An updated test indexer",
				"type": "private",
				"configured": true,
				"enabled": true,
				"fields": {
					"username": "newuser",
					"password": "newpass"
				}
			}`,
			wantError: false,
			wantConfig: &IndexerConfig{
				ID:          "test-indexer",
				Name:        "Updated Test Indexer",
				Description: "An updated test indexer",
				Type:        "private",
				Configured:  true,
				Enabled:     true,
				Fields: map[string]interface{}{
					"username": "newuser",
					"password": "newpass",
				},
			},
		},
		{
			name:      "validation error",
			indexerID: "test-indexer",
			config: &IndexerConfigRequest{
				Name: "", // Invalid empty name
			},
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"error": "Name is required"
			}`,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PUT" {
					t.Errorf("expected PUT request, got %s", r.Method)
				}
				expectedPath := "/api/v2.0/indexers/" + tt.indexerID + "/config"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}
				if r.URL.Query().Get("apikey") != testAPIKey {
					t.Errorf("expected API key %s, got %s", testAPIKey, r.URL.Query().Get("apikey"))
				}
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
				}

				// Verify request body
				var requestBody IndexerConfigRequest
				if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.responseStatus)
				w.Write([]byte(tt.responseBody))
			}))
			defer srv.Close()

			client, err := New(Settings{
				ApiURL: srv.URL,
				ApiKey: testAPIKey,
				Client: srv.Client(),
			})
			if err != nil {
				t.Fatal(err)
			}

			config, err := client.SaveIndexerConfig(context.Background(), tt.indexerID, tt.config)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if config == nil {
				t.Errorf("expected config, got nil")
				return
			}

			if config.ID != tt.wantConfig.ID {
				t.Errorf("expected ID %s, got %s", tt.wantConfig.ID, config.ID)
			}
			if config.Name != tt.wantConfig.Name {
				t.Errorf("expected Name %s, got %s", tt.wantConfig.Name, config.Name)
			}
		})
	}
}

func TestCreateIndexer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		config         *IndexerConfigRequest
		responseStatus int
		responseBody   string
		wantError      bool
		wantConfig     *IndexerConfig
	}{
		{
			name: "successful indexer creation",
			config: &IndexerConfigRequest{
				Name:        "New Test Indexer",
				Description: "A newly created test indexer",
				Type:        "public",
				Enabled:     true,
				Fields: map[string]interface{}{
					"url": "https://example.com",
				},
			},
			responseStatus: http.StatusCreated,
			responseBody: `{
				"id": "new-test-indexer",
				"name": "New Test Indexer",
				"description": "A newly created test indexer",
				"type": "public",
				"configured": true,
				"enabled": true,
				"fields": {
					"url": "https://example.com"
				}
			}`,
			wantError: false,
			wantConfig: &IndexerConfig{
				ID:          "new-test-indexer",
				Name:        "New Test Indexer",
				Description: "A newly created test indexer",
				Type:        "public",
				Configured:  true,
				Enabled:     true,
				Fields: map[string]interface{}{
					"url": "https://example.com",
				},
			},
		},
		{
			name: "creation with invalid data",
			config: &IndexerConfigRequest{
				Name: "", // Invalid empty name
			},
			responseStatus: http.StatusBadRequest,
			responseBody: `{
				"error": "Name is required"
			}`,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				if r.URL.Path != "/api/v2.0/indexers" {
					t.Errorf("expected path /api/v2.0/indexers, got %s", r.URL.Path)
				}
				if r.URL.Query().Get("apikey") != testAPIKey {
					t.Errorf("expected API key %s, got %s", testAPIKey, r.URL.Query().Get("apikey"))
				}
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
				}

				// Verify request body
				var requestBody IndexerConfigRequest
				if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.responseStatus)
				w.Write([]byte(tt.responseBody))
			}))
			defer srv.Close()

			client, err := New(Settings{
				ApiURL: srv.URL,
				ApiKey: testAPIKey,
				Client: srv.Client(),
			})
			if err != nil {
				t.Fatal(err)
			}

			config, err := client.CreateIndexer(context.Background(), tt.config)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if config == nil {
				t.Errorf("expected config, got nil")
				return
			}

			if config.ID != tt.wantConfig.ID {
				t.Errorf("expected ID %s, got %s", tt.wantConfig.ID, config.ID)
			}
			if config.Name != tt.wantConfig.Name {
				t.Errorf("expected Name %s, got %s", tt.wantConfig.Name, config.Name)
			}
		})
	}
}

func TestDeleteIndexer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		indexerID      string
		responseStatus int
		responseBody   string
		wantError      bool
	}{
		{
			name:           "successful indexer deletion",
			indexerID:      "test-indexer",
			responseStatus: http.StatusNoContent,
			responseBody:   "",
			wantError:      false,
		},
		{
			name:           "indexer not found",
			indexerID:      "nonexistent",
			responseStatus: http.StatusNotFound,
			responseBody:   `{"error": "Indexer not found"}`,
			wantError:      true,
		},
		{
			name:           "server error",
			indexerID:      "test-indexer",
			responseStatus: http.StatusInternalServerError,
			responseBody:   `{"error": "Internal server error"}`,
			wantError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("expected DELETE request, got %s", r.Method)
				}
				expectedPath := "/api/v2.0/indexers/" + tt.indexerID
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}
				if r.URL.Query().Get("apikey") != testAPIKey {
					t.Errorf("expected API key %s, got %s", testAPIKey, r.URL.Query().Get("apikey"))
				}

				w.WriteHeader(tt.responseStatus)
				if tt.responseBody != "" {
					w.Write([]byte(tt.responseBody))
				}
			}))
			defer srv.Close()

			client, err := New(Settings{
				ApiURL: srv.URL,
				ApiKey: testAPIKey,
				Client: srv.Client(),
			})
			if err != nil {
				t.Fatal(err)
			}

			err = client.DeleteIndexer(context.Background(), tt.indexerID)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestIndexerConfigValidation(t *testing.T) {
	t.Parallel()

	// Test that the JSON marshaling/unmarshaling works correctly
	config := &IndexerConfig{
		ID:          "test-id",
		Name:        "Test Indexer",
		Description: "A test indexer",
		Type:        "public",
		Configured:  true,
		Enabled:     true,
		Fields: map[string]interface{}{
			"username": "testuser",
			"password": "testpass",
			"enabled":  true,
			"timeout":  30,
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("failed to marshal config: %v", err)
	}

	// Unmarshal back
	var unmarshaledConfig IndexerConfig
	if err := json.Unmarshal(jsonData, &unmarshaledConfig); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}

	// Verify fields
	if unmarshaledConfig.ID != config.ID {
		t.Errorf("expected ID %s, got %s", config.ID, unmarshaledConfig.ID)
	}
	if unmarshaledConfig.Name != config.Name {
		t.Errorf("expected Name %s, got %s", config.Name, unmarshaledConfig.Name)
	}
	if unmarshaledConfig.Configured != config.Configured {
		t.Errorf("expected Configured %v, got %v", config.Configured, unmarshaledConfig.Configured)
	}
	if unmarshaledConfig.Enabled != config.Enabled {
		t.Errorf("expected Enabled %v, got %v", config.Enabled, unmarshaledConfig.Enabled)
	}

	// Verify fields map
	if len(unmarshaledConfig.Fields) != len(config.Fields) {
		t.Errorf("expected %d fields, got %d", len(config.Fields), len(unmarshaledConfig.Fields))
	}
	for key, expectedValue := range config.Fields {
		if actualValue, exists := unmarshaledConfig.Fields[key]; !exists {
			t.Errorf("expected field %s to exist", key)
		} else {
			// JSON unmarshaling converts numbers to float64, so we need to handle this
			switch expectedValue.(type) {
			case int:
				if actualFloat, ok := actualValue.(float64); !ok || int(actualFloat) != expectedValue {
					t.Errorf("expected field %s to be %v, got %v", key, expectedValue, actualValue)
				}
			default:
				if actualValue != expectedValue {
					t.Errorf("expected field %s to be %v, got %v", key, expectedValue, actualValue)
				}
			}
		}
	}
}
