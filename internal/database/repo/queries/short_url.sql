-- name: CreateShortUrl :one
INSERT INTO shortened_urls (
  original_url,
  short_url
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetByShortUrl :one
SELECT * FROM shortened_urls WHERE short_url = $1;

-- name: GetByOriginalUrl :one
SELECT * FROM shortened_urls WHERE original_url = $1;

-- name: DeleteByShortUrl :exec
DELETE FROM shortened_urls WHERE short_url = $1;

-- name: DeleteByOriginalUrl :exec
DELETE FROM shortened_urls WHERE original_url = $1;