-- name: CreateAppointment :one
INSERT INTO appointments (id, user_id, client_name, service_id, appointment_time, status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;