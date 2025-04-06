-- name: SaveURL :exec
INSERT INTO urls (short_url, original_url) VALUES ($1, $2);

-- name: GetURL :one
SELECT original_url FROM urls WHERE short_url = $1;

-- name: FindByOriginalURL :one
SELECT short_url FROM urls WHERE original_url = $1;