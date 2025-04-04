// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Appointment struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	ClientName      string    `json:"client_name"`
	ServiceID       uuid.UUID `json:"service_id"`
	AppointmentTime time.Time `json:"appointment_time"`
	Status          string    `json:"status"`
}

type Payment struct {
	ID         uuid.UUID      `json:"id"`
	ExternalID string         `json:"external_id"`
	Value      pgtype.Numeric `json:"value"`
	DueDate    time.Time      `json:"due_date"`
	Status     string         `json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type Schedule struct {
	ID              uuid.UUID   `json:"id"`
	UserID          uuid.UUID   `json:"user_id"`
	Weekday         int32       `json:"weekday"`
	StartTime       pgtype.Time `json:"start_time"`
	EndTime         pgtype.Time `json:"end_time"`
	IntervalMinutes int32       `json:"interval_minutes"`
}

type ScheduleOverride struct {
	ID         uuid.UUID   `json:"id"`
	ScheduleID uuid.UUID   `json:"schedule_id"`
	Date       time.Time   `json:"date"`
	StartTime  pgtype.Time `json:"start_time"`
	EndTime    pgtype.Time `json:"end_time"`
	Reason     pgtype.Text `json:"reason"`
}

type Service struct {
	ID          uuid.UUID      `json:"id"`
	TenantID    uuid.UUID      `json:"tenant_id"`
	Name        string         `json:"name"`
	Price       pgtype.Numeric `json:"price"`
	Description pgtype.Text    `json:"description"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type Tenant struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type UserService struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ServiceID uuid.UUID `json:"service_id"`
}
