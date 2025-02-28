package ports

import (
	"context"
	"github.com/google/uuid"
	"github.com/jpmoraess/pay/internal/domain"
)

type RegisterTenantInput struct {
	Name     string
	Email    string
	FullName string
	Password string
}

type RegisterTenantOutput struct {
	ID uuid.UUID
}

type TenantService interface {
	Register(ctx context.Context, input *RegisterTenantInput) (*RegisterTenantOutput, error)
}

type TenantRepository interface {
	Save(ctx context.Context, tenant *domain.Tenant) (*domain.Tenant, error)
}
