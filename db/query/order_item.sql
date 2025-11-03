-- name: CreateOrderItem :one
INSERT INTO order_item (
    order_id,
    menu_id,
    quantity,
    price
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetOrderItem :one
SELECT * FROM order_item
WHERE id = $1 LIMIT 1;

-- name: ListOrderItem :many
SELECT * FROM order_item
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateOrderItem :one
UPDATE order_item
SET order_id = $2,
    menu_id = $3,
    quantity = $4,
    note_item = $5,
    status = $6
WHERE id = $1
RETURNING *;

-- name: DeleteOrderItem :exec
DELETE FROM order_item
WHERE id = $1;



