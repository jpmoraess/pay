package ports

import (
	"context"
	"github.com/google/uuid"
	"github.com/jpmoraess/pay/internal/domain"
	"time"
)

type CreateSessionInput struct {
	ID           uuid.UUID
	Email        string
	RefreshToken string
	UserAgent    string
	ClientIP     string
	IsBlocked    bool
	ExpiresAt    time.Time
}

type CreateSessionOutput struct {
	ID           uuid.UUID
	Email        string
	RefreshToken string
	UserAgent    string
	ClientIP     string
	IsBlocked    bool
	ExpiresAt    time.Time
}

type GetSessionOutput struct {
	ID           uuid.UUID
	Email        string
	RefreshToken string
	UserAgent    string
	ClientIP     string
	IsBlocked    bool
	ExpiresAt    time.Time
}

type SessionService interface {
	Create(ctx context.Context, input *CreateSessionInput) (*CreateSessionOutput, error)
	GetSession(ctx context.Context, id uuid.UUID) (*GetSessionOutput, error)
}

type SessionRepository interface {
	Save(ctx context.Context, session *domain.Session) (*domain.Session, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Session, error)
}
