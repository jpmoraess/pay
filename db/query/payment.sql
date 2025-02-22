-- name: CreatePayment :one
INSERT INTO payments (id, external_id, value, due_date)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPayment :one
SELECT * FROM payments
WHERE id = $1
LIMIT 1;

-- name: GetPaymentByExternal :one
SELECT * FROM payments
WHERE external_id = $1
LIMIT 1;