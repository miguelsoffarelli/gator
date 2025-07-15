-- name: GetPostsForUser :many
SELECT * FROM posts
WHERE feed_id IN (
    SELECT id FROM feeds
    WHERE user_id = $1
)
ORDER BY created_at DESC
LIMIT $2;