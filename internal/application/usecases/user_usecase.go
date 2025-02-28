package usecases

import (
	"context"
	"errors"
	"github.com/jpmoraess/pay/config"
	db "github.com/jpmoraess/pay/db/sqlc"
	"github.com/jpmoraess/pay/internal/application/ports"
	"github.com/jpmoraess/pay/internal/domain"
	"github.com/jpmoraess/pay/token"
	"github.com/jpmoraess/pay/util"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type userUseCase struct {
	cfg            *config.Config
	tokenMaker     token.Maker
	repository     ports.UserRepository
	sessionService ports.SessionService
}

func NewUserUseCase(cfg *config.Config, tokenMaker token.Maker, repository ports.UserRepository, sessionService ports.SessionService) *userUseCase {
	return &userUseCase{cfg: cfg, tokenMaker: tokenMaker, repository: repository, sessionService: sessionService}
}

func (uc *userUseCase) Create(ctx context.Context, input *ports.CreateUserInput) (*ports.CreateUserOutput, error) {
	user, err := domain.NewUser(input.TenantID, input.FullName, input.Email, input.Password, input.Role)
	if err != nil {
		return nil, err
	}

	user, err = uc.repository.Save(ctx, user)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	output := &ports.CreateUserOutput{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	return output, nil
}

func (uc *userUseCase) Login(ctx context.Context, input *ports.LoginUserInput) (*ports.LoginUserOutput, error) {
	user, err := uc.repository.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, db.ErrNoRecordFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	err = util.CheckPassword(input.Password, user.Password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidPassword
		}
		return nil, err
	}

	accessToken, accessTokenPayload, err := uc.tokenMaker.CreateToken(user.TenantID, user.Email, user.Role, uc.cfg.AccessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenPayload, err := uc.tokenMaker.CreateToken(user.TenantID, user.Email, user.Role, uc.cfg.RefreshTokenDuration)
	if err != nil {
		return nil, err
	}

	userAgent, ok := ctx.Value("user-agent").(string)
	if !ok {
		userAgent = "unknown"
	}

	clientIP, ok := ctx.Value("client-ip").(string)
	if !ok {
		clientIP = "unknown"
	}

	session, err := uc.sessionService.Create(ctx, &ports.CreateSessionInput{
		ID:           refreshTokenPayload.ID,
		Email:        user.Email,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	})
	if err != nil {
		return nil, err
	}

	return &ports.LoginUserOutput{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
	}, nil
}
