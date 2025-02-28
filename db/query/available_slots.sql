-- name: GetAvailableSlots :many
-- Query to show available slots
WITH available_slots AS (
    -- Get time slots based on the service provider's schedule
    SELECT
        s.id AS schedule_id,
        s.user_id,
        s.weekday,
        generate_series(
            -- Generate the time slots within the available time range, using the dynamic date
                make_timestamp(EXTRACT(year FROM $2)::int, EXTRACT(month FROM $2)::int, EXTRACT(day FROM $2)::int,
                               EXTRACT(hour FROM s.start_time)::int, EXTRACT(minute FROM s.start_time)::int, 0),
                make_timestamp(EXTRACT(year FROM $2)::int, EXTRACT(month FROM $2)::int, EXTRACT(day FROM $2)::int,
                               EXTRACT(hour FROM s.end_time)::int, EXTRACT(minute FROM s.end_time)::int, 0),
                interval '1 minute' * s.interval_minutes
        ) AS available_time
    FROM
        schedules s
    WHERE
        s.user_id = $1 -- Service provider's ID (passed as a parameter)
),
     overridden_slots AS (
         -- Consider the slots that have been overridden (schedule changes)
         SELECT so.schedule_id, so.date, so.start_time, so.end_time
         FROM schedule_overrides so
         WHERE so.date = $2 -- Requested date (passed as a parameter)
     ),
     occupied_slots AS (
         -- Get the slots that are already occupied (appointments)
         SELECT a.appointment_time
         FROM appointments a
         WHERE a.user_id = $1
           AND a.appointment_time::date = $2 -- Requested date (passed as a parameter)
    )
-- Get available slots (not occupied or overridden)
SELECT
    available_slots.available_time
FROM
    available_slots
        LEFT JOIN
    overridden_slots
    ON available_slots.schedule_id = overridden_slots.schedule_id
        AND available_slots.available_time >= overridden_slots.start_time
        AND available_slots.available_time < overridden_slots.end_time
        LEFT JOIN
    occupied_slots
    ON available_slots.available_time = occupied_slots.appointment_time
WHERE
    available_slots.weekday = EXTRACT(dow FROM $2)::int -- Check if the weekday matches
    AND overridden_slots.schedule_id IS NULL -- Exclude overridden time slots
    AND occupied_slots.appointment_time IS NULL -- Exclude occupied time slots
ORDER BY available_slots.available_time;
