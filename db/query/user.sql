-- name: CreateUser :one
INSERT INTO users (
    username,
    hash_password,
    full_name,
    email
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: ListUser :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET 
    full_name = COALESCE(sqlc.narg(full_name), full_name),
    role = COALESCE(sqlc.narg(role), role),
    email = COALESCE(sqlc.narg(email), email)
WHERE id = $1
RETURNING *;

-- name: UpdateUserWithPassword :one
UPDATE users
SET hash_password = $2
WHERE id = $1
RETURNING *;


-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;



