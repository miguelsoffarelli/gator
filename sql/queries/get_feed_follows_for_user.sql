-- name: GetFeedFollowsForUser :many
WITH all_feed_follows AS (
    SELECT * FROM feed_follows
    WHERE feed_follows.user_id = $1
)
SELECT
    all_feed_follows.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM all_feed_follows
INNER JOIN feeds
    ON all_feed_follows.feed_id = feeds.id
INNER JOIN users
    ON all_feed_follows.user_id = users.id;