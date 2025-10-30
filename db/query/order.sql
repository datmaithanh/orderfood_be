-- name: CreateOrder :one
INSERT INTO orders (
    user_id,
    customer_id,
    table_id,
    status,
    total_price
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: ListOrder :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateOrder :one
UPDATE orders
SET user_id = $2,
    customer_id = $3,
    table_id = $4,
    status = $5,
    total_price = $6
WHERE id = $1
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;



