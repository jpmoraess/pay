package database

import (
	"context"
	"github.com/google/uuid"
	db "github.com/jpmoraess/pay/db/sqlc"
	"github.com/jpmoraess/pay/internal/domain"
)

type sessionRepository struct {
	db db.Store
}

func NewSessionRepository(db db.Store) *sessionRepository {
	return &sessionRepository{db: db}
}

func newSession(sessionDB db.Session) *domain.Session {
	return &domain.Session{
		ID:           sessionDB.ID,
		Email:        sessionDB.Email,
		RefreshToken: sessionDB.RefreshToken,
		UserAgent:    sessionDB.UserAgent,
		ClientIP:     sessionDB.ClientIp,
		IsBlocked:    sessionDB.IsBlocked,
		ExpiresAt:    sessionDB.ExpiresAt,
	}
}

func newCreateSessionParams(session *domain.Session) db.CreateSessionParams {
	return db.CreateSessionParams{
		ID:           session.ID,
		Email:        session.Email,
		RefreshToken: session.RefreshToken,
		UserAgent:    session.UserAgent,
		ClientIp:     session.ClientIP,
		IsBlocked:    session.IsBlocked,
		ExpiresAt:    session.ExpiresAt,
	}
}

func (repo *sessionRepository) Save(ctx context.Context, session *domain.Session) (*domain.Session, error) {
	sessionDB, err := repo.db.CreateSession(ctx, newCreateSessionParams(session))
	if err != nil {
		return nil, err
	}
	return newSession(sessionDB), nil
}

func (repo *sessionRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	sessionDB, err := repo.db.GetSession(ctx, id)
	if err != nil {
		return nil, err
	}
	return newSession(sessionDB), nil
}
