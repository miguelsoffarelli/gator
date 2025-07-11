-- name: GetUserByName :one
SELECT * FROM users
WHERE name = $1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;