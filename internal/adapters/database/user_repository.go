package database

import (
	"context"
	db "github.com/jpmoraess/pay/db/sqlc"
	"github.com/jpmoraess/pay/internal/domain"
)

type userRepository struct {
	db db.Store
}

func NewUserRepository(db db.Store) *userRepository {
	return &userRepository{db: db}
}

func newUser(userDB db.User) *domain.User {
	return &domain.User{
		ID:       userDB.ID,
		TenantID: userDB.TenantID,
		Email:    userDB.Email,
		Password: userDB.Password,
		FullName: userDB.FullName,
		Role:     userDB.Role,
	}
}

func newCreateUserParams(user *domain.User) db.CreateUserParams {
	return db.CreateUserParams{
		ID:       user.ID,
		TenantID: user.TenantID,
		Email:    user.Email,
		Password: user.Password,
		FullName: user.FullName,
		Role:     user.Role,
	}
}

func (repo *userRepository) Save(ctx context.Context, user *domain.User) (*domain.User, error) {
	userDB, err := repo.db.CreateUser(ctx, newCreateUserParams(user))
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, err
		}
		return nil, err
	}

	return newUser(userDB), nil
}

func (repo *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	userDB, err := repo.db.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return newUser(userDB), nil
}
