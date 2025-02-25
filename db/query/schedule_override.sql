-- name: CreateScheduleOverride :one
INSERT INTO schedule_overrides (id, schedule_id, date, start_time, end_time, reason)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;