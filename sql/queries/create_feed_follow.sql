-- name: CreateFeedFollow :one
WITH feed_follow_insert AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT
    feed_follow_insert.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM feed_follow_insert
INNER JOIN feeds
ON feed_follow_insert.feed_id = feeds.id
INNER JOIN users
ON feed_follow_insert.user_id = users.id;
