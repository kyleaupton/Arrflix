package repo

import (
	"context"
)

type SetupRepo interface {
	GetSystemInitialized(ctx context.Context) (bool, error)
	SetSystemInitialized(ctx context.Context) error
	CountUsers(ctx context.Context) (int64, error)
}

func (r *Repository) GetSystemInitialized(ctx context.Context) (bool, error) {
	return r.Q.GetSystemInitialized(ctx)
}

func (r *Repository) SetSystemInitialized(ctx context.Context) error {
	return r.Q.SetSystemInitialized(ctx)
}

func (r *Repository) CountUsers(ctx context.Context) (int64, error) {
	return r.Q.CountUsers(ctx)
}
