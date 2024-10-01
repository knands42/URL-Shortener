-- name: CreateShortUrl :one
INSERT INTO shortened_urls (
  original_url,
  short_url
) VALUES (
  $1, $2
) RETURNING *;
