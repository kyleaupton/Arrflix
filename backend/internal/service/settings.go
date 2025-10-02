package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/kyleaupton/snaggle/backend/internal/repo"
)

type SettingsService struct {
	repo *repo.Repository
	mu   sync.RWMutex
	mem  map[string]any
}

func NewSettingsService(r *repo.Repository) *SettingsService {
	return &SettingsService{repo: r, mem: make(map[string]any)}
}

// GetAll returns a materialized map of settings with defaults applied and caches it.
func (s *SettingsService) GetAll(ctx context.Context) (map[string]any, error) {
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make(map[string]any, len(rows))
	for _, it := range rows {
		switch it.Type {
		case string(SettingText):
			var v string
			_ = json.Unmarshal(it.ValueJson, &v)
			out[it.Key] = v
		case string(SettingBool):
			var v bool
			_ = json.Unmarshal(it.ValueJson, &v)
			out[it.Key] = v
		case string(SettingInt):
			var v int64
			_ = json.Unmarshal(it.ValueJson, &v)
			out[it.Key] = v
		case string(SettingJSON):
			var v any
			_ = json.Unmarshal(it.ValueJson, &v)
			out[it.Key] = v
		default:
			var v any
			_ = json.Unmarshal(it.ValueJson, &v)
			out[it.Key] = v
		}
	}
	for k, spec := range Registry {
		if _, ok := out[k]; !ok {
			out[k] = spec.Default
		}
	}
	s.mu.Lock()
	s.mem = out
	s.mu.Unlock()
	return out, nil
}

// Set validates and persists a single key/value according to the registry.
func (s *SettingsService) Set(ctx context.Context, key string, val any) error {
	spec, ok := Registry[key]
	if !ok {
		return fmt.Errorf("unknown setting key")
	}
	var (
		typ = string(spec.Type)
		b   []byte
		err error
	)
	switch spec.Type {
	case SettingText:
		if _, ok := val.(string); !ok {
			return fmt.Errorf("invalid type; want text")
		}
	case SettingBool:
		if _, ok := val.(bool); !ok {
			return fmt.Errorf("invalid type; want bool")
		}
	case SettingInt:
		switch t := val.(type) {
		case float64:
			val = int64(t)
		case int:
			val = int64(t)
		case int64:
		default:
			return fmt.Errorf("invalid type; want int")
		}
	case SettingJSON:
		// any
	}
	if b, err = json.Marshal(val); err != nil {
		return err
	}
	if err := s.repo.Upsert(ctx, key, typ, b); err != nil {
		return err
	}
	s.mu.Lock()
	s.mem[key] = val
	s.mu.Unlock()
	return nil
}
