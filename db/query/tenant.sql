-- name: CreateTenant :one
INSERT INTO tenants (id, name, email, password)
VALUES ($1, $2, $3, $4)
RETURNING *;