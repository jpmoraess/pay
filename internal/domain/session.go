package domain

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID           uuid.UUID
	Email        string
	RefreshToken string
	UserAgent    string
	ClientIP     string
	IsBlocked    bool
	ExpiresAt    time.Time
}

func NewSession(id uuid.UUID, email, refreshToken, userAgent, clientIP string, isBlocked bool, expiresAt time.Time) (*Session, error) {
	return &Session{
		ID:           id,
		Email:        email,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		IsBlocked:    isBlocked,
		ExpiresAt:    expiresAt,
	}, nil
}
