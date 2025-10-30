-- name: CreatePayment :one
INSERT INTO payments (
    order_id,
    amount,
    payment_method,
    status
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetPayment :one
SELECT * FROM payments
WHERE id = $1 LIMIT 1;

-- name: ListPayment :many
SELECT * FROM payments
WHERE order_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdatePayment :one
UPDATE payments
SET order_id = $2,
    amount = $3,
    payment_method = $4,
    status = $5
WHERE id = $1
RETURNING *;

-- name: DeletePayment :exec
DELETE FROM payments
WHERE id = $1;



