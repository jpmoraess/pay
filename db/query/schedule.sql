-- name: CreateSchedule :one
INSERT INTO schedules (user_id, weekday, start_time, end_time, interval_minutes)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSchedule :one
SELECT * FROM schedules
WHERE id = $1
LIMIT 1;

