-- name: CreatePayment :one
INSERT INTO payments (
    order_id,
    amount,
    payment_method
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetPayment :one
SELECT * FROM payments
WHERE id = $1 LIMIT 1;

-- name: ListPayment :many
SELECT * FROM payments
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdatePaymentStatus :one
UPDATE payments
SET status = $2
WHERE id = $1
RETURNING *;


-- name: DeletePayment :exec
DELETE FROM payments
WHERE id = $1;



