package repo

import (
	"context"

	dbgen "github.com/kyleaupton/arrflix/internal/db/sqlc"
)

type SettingsRepo interface {
	Get(ctx context.Context, key string) (dbgen.AppSetting, error)
	List(ctx context.Context) ([]dbgen.AppSetting, error)
	Upsert(ctx context.Context, key, typ string, valueJson []byte) error
}

func (r *Repository) Get(ctx context.Context, key string) (dbgen.AppSetting, error) {
	return r.Q.GetSetting(ctx, key)
}
func (r *Repository) List(ctx context.Context) ([]dbgen.AppSetting, error) {
	return r.Q.ListSettings(ctx)
}
func (r *Repository) Upsert(ctx context.Context, key, typ string, valueJson []byte) error {
	return r.Q.UpsertSetting(ctx, dbgen.UpsertSettingParams{Key: key, Type: typ, ValueJson: valueJson})
}
