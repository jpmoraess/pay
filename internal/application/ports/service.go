package ports

import (
	"context"
	"github.com/google/uuid"
	"github.com/jpmoraess/pay/internal/domain"
)

type CreateServiceInput struct {
	TenantID    uuid.UUID
	Name        string
	Price       float64
	Description string
}

type CreateServiceOutput struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Name        string
	Price       float64
	Description string
}

type ServiceService interface {
	Create(ctx context.Context, input *CreateServiceInput) (*CreateServiceOutput, error)
}

type ServiceRepository interface {
	Save(ctx context.Context, service *domain.Service) (*domain.Service, error)
}
