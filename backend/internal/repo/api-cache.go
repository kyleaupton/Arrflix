package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/snaggle/backend/internal/db/sqlc"
)

type CacheEntry struct {
	Key         string
	Category    *string
	Response    []byte
	Status      int32
	ContentType *string
	Headers     []byte
	StoredAt    pgtype.Timestamptz
	ExpiresAt   pgtype.Timestamptz
}

func (r *Repository) GetApiCache(ctx context.Context, key string) (CacheEntry, bool, error) {
	row, err := r.Q.GetApiCache(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// cache miss: no rows
			return CacheEntry{}, false, nil
		}

		// real error
		return CacheEntry{}, false, err
	}

	if row.ExpiresAt.Time.Before(time.Now()) {
		return CacheEntry{}, false, nil
	}

	return CacheEntry{
		Key:         row.Key,
		Category:    row.Category,
		Response:    row.Response,
		Status:      row.Status,
		ContentType: row.ContentType,
		Headers:     row.Headers,
		StoredAt:    row.StoredAt,
		ExpiresAt:   row.ExpiresAt,
	}, true, nil
}

func (r *Repository) UpsertApiCache(ctx context.Context, key string, category *string, response []byte, status int, contentType *string, headers []byte, ttl time.Duration) error {
	expires := pgtype.Timestamptz{}
	_ = expires.Scan(time.Now().Add(ttl))

	fmt.Println("key", key)
	fmt.Println("expires", expires)

	return r.Q.UpsertApiCache(ctx, dbgen.UpsertApiCacheParams{
		Key:         key,
		Category:    category,
		Response:    response,
		Status:      int32(status),
		ContentType: contentType,
		Headers:     headers,
		ExpiresAt:   expires,
	})
}

func (r *Repository) DeleteExpiredApiCache(ctx context.Context) error {
	return r.Q.DeleteExpiredApiCache(ctx)
}
