package domain

import (
	"github.com/google/uuid"
	"github.com/jpmoraess/pay/util"
	"time"
)

type User struct {
	ID        uuid.UUID
	FullName  string
	Email     string
	Password  string
	CreatedAt time.Time
}

func NewUser(fullName, email, password string) (*User, error) {
	hashedPwd, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        uuid.New(),
		FullName:  fullName,
		Email:     email,
		Password:  hashedPwd,
		CreatedAt: time.Now(),
	}, nil
}
