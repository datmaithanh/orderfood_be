-- name: CreateTable :one
INSERT INTO tables (
    name,
    qr_text,
    qr_image_url
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetMaxTableID :one
SELECT COALESCE(MAX(id), 0) FROM tables;

-- name: GetTable :one
SELECT * FROM tables
WHERE id = $1 LIMIT 1;

-- name: ListTable :many
SELECT * FROM tables
ORDER BY id
LIMIT $1
OFFSET $2;

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



