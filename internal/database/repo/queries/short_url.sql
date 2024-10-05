-- name: CreateHash :one
INSERT INTO shortened_urls (
  original_url,
  hash
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetByHash :one
SELECT * FROM shortened_urls WHERE hash = $1;

-- name: GetByOriginalUrl :one
SELECT * FROM shortened_urls WHERE original_url = $1;

-- name: DeleteByHash :exec
DELETE FROM shortened_urls WHERE hash = $1;

-- name: DeleteByOriginalUrl :exec
DELETE FROM shortened_urls WHERE original_url = $1;

-- name: IncreaseNumberOfAccesses :exec
UPDATE shortened_urls SET number_of_accesses = number_of_accesses + 1 WHERE hash = $1;