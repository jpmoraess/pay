-- name: GetSchedules :many
SELECT 
    schedules.start_time, 
    schedules.end_time, 
    schedules.interval_minutes, 
    COALESCE(overrides.available, TRUE) as available
FROM schedules
LEFT JOIN schedule_overrides overrides
    ON s.user_id = overrides.user_id
    AND overrides.date = $1
WHERE s.user_id = $2
AND s.weekday = EXTRACT(DOW FROM $1)
ORDER BY schedules.start_time;