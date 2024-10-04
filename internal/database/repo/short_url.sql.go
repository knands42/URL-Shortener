// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: short_url.sql

package repo

import (
	"context"
)

const createHash = `-- name: CreateHash :one
INSERT INTO shortened_urls (
  original_url,
  hash
) VALUES (
  $1, $2
) RETURNING id, original_url, hash, number_of_access, created_at, updated_at
`

type CreateHashParams struct {
	OriginalUrl string `json:"original_url"`
	Hash        string `json:"hash"`
}

func (q *Queries) CreateHash(ctx context.Context, arg CreateHashParams) (ShortenedUrl, error) {
	row := q.db.QueryRow(ctx, createHash, arg.OriginalUrl, arg.Hash)
	var i ShortenedUrl
	err := row.Scan(
		&i.ID,
		&i.OriginalUrl,
		&i.Hash,
		&i.NumberOfAccess,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteByHash = `-- name: DeleteByHash :exec
DELETE FROM shortened_urls WHERE hash = $1
`

func (q *Queries) DeleteByHash(ctx context.Context, hash string) error {
	_, err := q.db.Exec(ctx, deleteByHash, hash)
	return err
}

const deleteByOriginalUrl = `-- name: DeleteByOriginalUrl :exec
DELETE FROM shortened_urls WHERE original_url = $1
`

func (q *Queries) DeleteByOriginalUrl(ctx context.Context, originalUrl string) error {
	_, err := q.db.Exec(ctx, deleteByOriginalUrl, originalUrl)
	return err
}

const getByHash = `-- name: GetByHash :one
SELECT id, original_url, hash, number_of_access, created_at, updated_at FROM shortened_urls WHERE hash = $1
`

func (q *Queries) GetByHash(ctx context.Context, hash string) (ShortenedUrl, error) {
	row := q.db.QueryRow(ctx, getByHash, hash)
	var i ShortenedUrl
	err := row.Scan(
		&i.ID,
		&i.OriginalUrl,
		&i.Hash,
		&i.NumberOfAccess,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getByOriginalUrl = `-- name: GetByOriginalUrl :one
SELECT id, original_url, hash, number_of_access, created_at, updated_at FROM shortened_urls WHERE original_url = $1
`

func (q *Queries) GetByOriginalUrl(ctx context.Context, originalUrl string) (ShortenedUrl, error) {
	row := q.db.QueryRow(ctx, getByOriginalUrl, originalUrl)
	var i ShortenedUrl
	err := row.Scan(
		&i.ID,
		&i.OriginalUrl,
		&i.Hash,
		&i.NumberOfAccess,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
