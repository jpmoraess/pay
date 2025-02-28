package usecases

import (
	"context"
	"errors"
	db "github.com/jpmoraess/pay/db/sqlc"
	"github.com/jpmoraess/pay/internal/application/ports"
	"github.com/jpmoraess/pay/internal/domain"
)

var (
	ErrTenantAlreadyExists = errors.New("tenant already exists")
)

type tenantUseCase struct {
	repository ports.TenantRepository
}

func NewTenantUseCase(repository ports.TenantRepository) *tenantUseCase {
	return &tenantUseCase{
		repository: repository,
	}
}

func (uc *tenantUseCase) Register(ctx context.Context, input *ports.RegisterTenantInput) (*ports.RegisterTenantOutput, error) {
	tenant, err := domain.NewTenant(input.Email, input.Name, input.FullName, input.Password)
	if err != nil {
		return nil, err
	}

	tenant, err = uc.repository.Save(ctx, tenant)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, ErrTenantAlreadyExists
		}
		return nil, err
	}

	output := &ports.RegisterTenantOutput{
		ID: tenant.ID,
	}

	return output, nil
}
