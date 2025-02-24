package database

import (
	"context"
	"github.com/hibiken/asynq"
	db "github.com/jpmoraess/pay/db/sqlc"
	"github.com/jpmoraess/pay/internal/domain"
	"github.com/jpmoraess/pay/worker"
	"time"
)

type userRepository struct {
	db              db.Store
	taskDistributor worker.TaskDistributor
}

func NewUserRepository(db db.Store, taskDistributor worker.TaskDistributor) *userRepository {
	return &userRepository{db: db, taskDistributor: taskDistributor}
}

func newUser(userDB db.User) *domain.User {
	return &domain.User{
		ID:       userDB.ID,
		Email:    userDB.Email,
		Password: userDB.Password,
		FullName: userDB.FullName,
	}
}

func (u *userRepository) newCreateUserParams(ctx context.Context, user *domain.User) db.CreateUserTxParams {
	return db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			ID:       user.ID,
			Email:    user.Email,
			Password: user.Password,
			FullName: user.FullName,
			Role:     "simple",
		},
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Email: user.Email,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(5 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			return u.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}
}

func (repo *userRepository) Save(ctx context.Context, user *domain.User) (*domain.User, error) {
	userDB, err := repo.db.CreateUserTx(ctx, repo.newCreateUserParams(ctx, user))
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, err
		}
		return nil, err
	}

	return newUser(userDB.User), nil
}

func (repo *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	userDB, err := repo.db.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return newUser(userDB), nil
}
