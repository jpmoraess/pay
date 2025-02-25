package usecases

import (
	"context"
	"github.com/jpmoraess/pay/internal/application/ports"
	"github.com/jpmoraess/pay/internal/domain"
)

type serviceUseCase struct {
	repository ports.ServiceRepository
}

func NewServiceUseCase(repository ports.ServiceRepository) *serviceUseCase {
	return &serviceUseCase{
		repository: repository,
	}
}

func (uc *serviceUseCase) Create(ctx context.Context, input *ports.CreateServiceInput) (*ports.CreateServiceOutput, error) {
	service, err := domain.NewService(input.TenantID, input.Name, input.Price, input.Description)
	if err != nil {
		return nil, err
	}

	service, err = uc.repository.Save(ctx, service)
	if err != nil {
		return nil, err
	}

	output := &ports.CreateServiceOutput{
		ID:          service.ID,
		TenantID:    service.TenantID,
		Name:        service.Name,
		Price:       service.Price,
		Description: service.Description,
	}

	return output, nil
}
