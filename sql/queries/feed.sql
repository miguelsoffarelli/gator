-- name: GetFeed :one
SELECT * FROM feeds
WHERE name = $1;