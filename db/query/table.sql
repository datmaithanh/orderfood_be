-- name: CreateTable :one
INSERT INTO tables (
    name,
    qr_text,
    qr_image_url,
    status
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetTable :one
SELECT * FROM tables
WHERE id = $1 LIMIT 1;

-- name: ListTable :many
SELECT * FROM tables
WHERE name = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateTable :one
UPDATE tables
SET name = $2,
    qr_text = $3,
    qr_image_url = $4,
    status = $5
WHERE id = $1
RETURNING *;

-- name: DeleteTable :exec
DELETE FROM tables
WHERE id = $1;



