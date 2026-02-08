package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/arrflix/internal/db/sqlc"
)

func (r *Repository) CreateInvite(ctx context.Context, email string, invitedBy pgtype.UUID) (dbgen.UserInvite, error) {
	return r.Q.CreateInvite(ctx, dbgen.CreateInviteParams{
		Email:     email,
		InvitedBy: invitedBy,
	})
}

func (r *Repository) GetInviteByEmail(ctx context.Context, email string) (dbgen.UserInvite, error) {
	return r.Q.GetInviteByEmail(ctx, email)
}

func (r *Repository) ClaimInvite(ctx context.Context, id pgtype.UUID) error {
	return r.Q.ClaimInvite(ctx, id)
}

func (r *Repository) ListInvites(ctx context.Context) ([]dbgen.UserInvite, error) {
	return r.Q.ListInvites(ctx)
}

func (r *Repository) DeleteInvite(ctx context.Context, id pgtype.UUID) error {
	return r.Q.DeleteInvite(ctx, id)
}
