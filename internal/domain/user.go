package domain

import (
	"github.com/google/uuid"
	"github.com/jpmoraess/pay/util"
	"time"
)

type User struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	FullName  string
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
}

func NewUser(tenantID uuid.UUID, fullName, email, password, role string) (*User, error) {
	hashedPwd, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        uuid.New(),
		TenantID:  tenantID,
		FullName:  fullName,
		Email:     email,
		Password:  hashedPwd,
		Role:      role,
		CreatedAt: time.Now(),
	}, nil
}
