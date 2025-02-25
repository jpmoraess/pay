-- name: CreateService :one
INSERT INTO services (id, tenant_id, name, price, description)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetService :one
SELECT * FROM services
WHERE id = $1
LIMIT 1;

-- name: ListServices :many
SELECT * FROM services
WHERE tenant_id = $1;