package domain

import "github.com/google/uuid"

type Service struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Name        string
	Price       float64
	Description string
}

func NewService(tenantID uuid.UUID, name string, price float64, description string) (*Service, error) {
	return &Service{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Name:        name,
		Price:       price,
		Description: description,
	}, nil
}
