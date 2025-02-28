package domain

import (
	"github.com/google/uuid"
	"github.com/jpmoraess/pay/util"
)

type Tenant struct {
	ID       uuid.UUID
	Name     string
	FullName string
	Email    string
	Password string
}

func NewTenant(email string, name string, fullName string, password string) (*Tenant, error) {
	hashedPwd, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}

	return &Tenant{
		ID:       uuid.New(),
		Name:     name,
		FullName: fullName,
		Email:    email,
		Password: hashedPwd,
	}, nil
}
