-- name: CreateOrderItem :one
INSERT INTO order_item (
    order_id,
    menu_id,
    quantity,
    price,
    note_item,
    status
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetOrderItem :one
SELECT * FROM order_item
WHERE id = $1 LIMIT 1;

-- name: ListOrderItem :many
SELECT * FROM order_item
WHERE order_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateOrderItem :one
UPDATE order_item
SET order_id = $2,
    menu_id = $3,
    quantity = $4,
    price = $5,
    note_item = $6,
    status = $7
WHERE id = $1
RETURNING *;

-- name: DeleteOrderItem :exec
DELETE FROM order_item
WHERE id = $1;



