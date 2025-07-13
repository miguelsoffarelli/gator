-- name: GetFeedByName :one
SELECT * FROM feeds
WHERE name = $1;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1;