package ports

import (
	"context"
	"github.com/google/uuid"
	"github.com/jpmoraess/pay/internal/domain"
	"time"
)

type CreateUserInput struct {
	FullName string
	Email    string
	Password string
}

type CreateUserOutput struct {
	FullName  string
	Email     string
	CreatedAt time.Time
}

type LoginUserInput struct {
	Email    string
	Password string
}

type LoginUserOutput struct {
	SessionID             uuid.UUID
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

type UserService interface {
	Create(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error)
	Login(ctx context.Context, input *LoginUserInput) (*LoginUserOutput, error)
}

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}
