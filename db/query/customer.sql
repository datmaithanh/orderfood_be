-- name: CreateCustomer :one
INSERT INTO customers (
    full_name,
    phone_number,
    email
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetCustomer :one
SELECT * FROM customers
WHERE id = $1 LIMIT 1;

-- name: ListCustomer :many
SELECT * FROM customers
WHERE phone_number = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: DeleteCustomer :exec
DELETE FROM customers
WHERE id = $1;
