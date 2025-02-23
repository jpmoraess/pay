package usecases

import (
	"context"
	"errors"
	"github.com/hibiken/asynq"
	"github.com/jpmoraess/pay/config"
	db "github.com/jpmoraess/pay/db/sqlc"
	"github.com/jpmoraess/pay/internal/application/ports"
	"github.com/jpmoraess/pay/internal/domain"
	"github.com/jpmoraess/pay/token"
	"github.com/jpmoraess/pay/util"
	"github.com/jpmoraess/pay/worker"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrSendEmailVerifyFailed = errors.New("failed to distribute task to send verify email")
)

type userUseCase struct {
	cfg             *config.Config
	tokenMaker      token.Maker
	repository      ports.UserRepository
	sessionService  ports.SessionService
	taskDistributor worker.TaskDistributor
}

func NewUserUseCase(
	cfg *config.Config,
	tokenMaker token.Maker,
	repository ports.UserRepository,
	sessionService ports.SessionService,
	taskDistributor worker.TaskDistributor,
) *userUseCase {
	return &userUseCase{
		cfg:             cfg,
		tokenMaker:      tokenMaker,
		repository:      repository,
		sessionService:  sessionService,
		taskDistributor: taskDistributor,
	}
}

func (uc *userUseCase) Create(ctx context.Context, input *ports.CreateUserInput) (*ports.CreateUserOutput, error) {
	user, err := domain.NewUser(input.FullName, input.Email, input.Password)
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

	taskPayload := &worker.PayloadSendVerifyEmail{
		Email: user.Email,
	}
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(5 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}
	err = uc.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
	if err != nil {
		return nil, ErrSendEmailVerifyFailed
	}

	output := &ports.CreateUserOutput{
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

	accessToken, accessTokenPayload, err := uc.tokenMaker.CreateToken(user.Email, "simple", uc.cfg.AccessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenPayload, err := uc.tokenMaker.CreateToken(user.Email, "simple", uc.cfg.RefreshTokenDuration)
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
