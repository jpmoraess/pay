package usecases

import (
	"context"
	"errors"
	"github.com/google/uuid"
	db "github.com/jpmoraess/pay/db/sqlc"
	"github.com/jpmoraess/pay/internal/application/ports"
	"github.com/jpmoraess/pay/internal/domain"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

type sessionUseCase struct {
	repository ports.SessionRepository
}

func NewSessionUseCase(repository ports.SessionRepository) *sessionUseCase {
	return &sessionUseCase{repository: repository}
}

func (uc *sessionUseCase) Create(ctx context.Context, input *ports.CreateSessionInput) (*ports.CreateSessionOutput, error) {
	session, err := domain.NewSession(input.ID, input.Email, input.RefreshToken, input.UserAgent, input.ClientIP, input.IsBlocked, input.ExpiresAt)
	if err != nil {
		return nil, err
	}

	session, err = uc.repository.Save(ctx, session)
	if err != nil {
		return nil, err
	}

	output := &ports.CreateSessionOutput{
		ID:           session.ID,
		Email:        session.Email,
		RefreshToken: session.RefreshToken,
		UserAgent:    session.UserAgent,
		ClientIP:     session.ClientIP,
		IsBlocked:    session.IsBlocked,
		ExpiresAt:    session.ExpiresAt,
	}

	return output, nil
}

func (uc *sessionUseCase) GetSession(ctx context.Context, id uuid.UUID) (*ports.GetSessionOutput, error) {
	session, err := uc.repository.Get(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRecordFound) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	output := &ports.GetSessionOutput{
		ID:           session.ID,
		Email:        session.Email,
		RefreshToken: session.RefreshToken,
		UserAgent:    session.UserAgent,
		ClientIP:     session.ClientIP,
		IsBlocked:    session.IsBlocked,
		ExpiresAt:    session.ExpiresAt,
	}

	return output, nil
}
