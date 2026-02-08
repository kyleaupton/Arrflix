package service

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/arrflix/internal/db/sqlc"
	"github.com/kyleaupton/arrflix/internal/repo"
)

type InvitesService struct {
	repo *repo.Repository
}

func NewInvitesService(r *repo.Repository) *InvitesService {
	return &InvitesService{repo: r}
}

// Create adds a new email to the invite list.
func (s *InvitesService) Create(ctx context.Context, email string, invitedBy pgtype.UUID) (dbgen.UserInvite, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return dbgen.UserInvite{}, errors.New("email required")
	}
	invite, err := s.repo.CreateInvite(ctx, email, invitedBy)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return dbgen.UserInvite{}, errors.New("email already invited")
		}
		return dbgen.UserInvite{}, err
	}
	return invite, nil
}

// List returns all invites.
func (s *InvitesService) List(ctx context.Context) ([]dbgen.UserInvite, error) {
	return s.repo.ListInvites(ctx)
}

// Delete removes an invite.
func (s *InvitesService) Delete(ctx context.Context, id pgtype.UUID) error {
	return s.repo.DeleteInvite(ctx, id)
}

// CheckAndClaim looks up an unclaimed invite by email and marks it as claimed.
// Returns nil if the invite exists and was claimed, an error otherwise.
func (s *InvitesService) CheckAndClaim(ctx context.Context, email string) error {
	invite, err := s.repo.GetInviteByEmail(ctx, email)
	if err != nil {
		return errors.New("no invite found for this email")
	}
	return s.repo.ClaimInvite(ctx, invite.ID)
}
