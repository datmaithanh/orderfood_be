-- name: CreateMenu :one
INSERT INTO menus (
    name,
    price,
    category_id
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetMenu :one
SELECT * FROM menus
WHERE id = $1 LIMIT 1;

-- name: ListMenu :many
SELECT * FROM menus
WHERE name = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateMenu :one
UPDATE menus
SET name = $2,
    price = $3,
    category_id = $4
WHERE id = $1
RETURNING *;

-- name: DeleteMenu :exec
DELETE FROM menus
WHERE id = $1;



