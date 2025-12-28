package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/snaggle/backend/internal/db/sqlc"
)

type NameTemplateRepo interface {
	ListNameTemplates(ctx context.Context) ([]dbgen.NameTemplate, error)
	GetNameTemplate(ctx context.Context, id pgtype.UUID) (dbgen.NameTemplate, error)
	GetDefaultNameTemplate(ctx context.Context, typ string) (dbgen.NameTemplate, error)
	CreateNameTemplate(ctx context.Context, name, typ, template string, showTemplate, seasonTemplate *string, isDefault bool) (dbgen.NameTemplate, error)
	UpdateNameTemplate(ctx context.Context, id pgtype.UUID, name, typ, template string, showTemplate, seasonTemplate *string, isDefault bool) (dbgen.NameTemplate, error)
	DeleteNameTemplate(ctx context.Context, id pgtype.UUID) error
}

func (r *Repository) ListNameTemplates(ctx context.Context) ([]dbgen.NameTemplate, error) {
	return r.Q.ListNameTemplates(ctx)
}

func (r *Repository) GetNameTemplate(ctx context.Context, id pgtype.UUID) (dbgen.NameTemplate, error) {
	return r.Q.GetNameTemplate(ctx, id)
}

func (r *Repository) GetDefaultNameTemplate(ctx context.Context, typ string) (dbgen.NameTemplate, error) {
	return r.Q.GetDefaultNameTemplate(ctx, typ)
}

func (r *Repository) CreateNameTemplate(ctx context.Context, name, typ, template string, showTemplate, seasonTemplate *string, isDefault bool) (dbgen.NameTemplate, error) {
	return r.Q.CreateNameTemplate(ctx, dbgen.CreateNameTemplateParams{
		Name:                 name,
		Type:                 typ,
		Template:             template,
		SeriesShowTemplate:   showTemplate,
		SeriesSeasonTemplate: seasonTemplate,
		IsDefault:            isDefault,
	})
}

func (r *Repository) UpdateNameTemplate(ctx context.Context, id pgtype.UUID, name, typ, template string, showTemplate, seasonTemplate *string, isDefault bool) (dbgen.NameTemplate, error) {
	return r.Q.UpdateNameTemplate(ctx, dbgen.UpdateNameTemplateParams{
		ID:                   id,
		Name:                 name,
		Type:                 typ,
		Template:             template,
		SeriesShowTemplate:   showTemplate,
		SeriesSeasonTemplate: seasonTemplate,
		IsDefault:            isDefault,
	})
}

func (r *Repository) DeleteNameTemplate(ctx context.Context, id pgtype.UUID) error {
	return r.Q.DeleteNameTemplate(ctx, id)
}
