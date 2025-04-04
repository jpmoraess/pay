// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: available_slots.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const getAvailableSlots = `-- name: GetAvailableSlots :many
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
ORDER BY available_slots.available_time
`

type GetAvailableSlotsParams struct {
	UserID  uuid.UUID `json:"user_id"`
	Extract time.Time `json:"extract"`
}

// Query to show available slots
// Get available slots (not occupied or overridden)
func (q *Queries) GetAvailableSlots(ctx context.Context, arg GetAvailableSlotsParams) ([]int64, error) {
	rows, err := q.db.Query(ctx, getAvailableSlots, arg.UserID, arg.Extract)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int64{}
	for rows.Next() {
		var available_time int64
		if err := rows.Scan(&available_time); err != nil {
			return nil, err
		}
		items = append(items, available_time)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
