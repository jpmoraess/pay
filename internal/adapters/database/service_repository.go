package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/jpmoraess/pay/db/sqlc"
	"github.com/jpmoraess/pay/internal/domain"
	"github.com/shopspring/decimal"
)

type serviceRepository struct {
	db db.Store
}

func NewServiceRepository(db db.Store) *serviceRepository {
	return &serviceRepository{db: db}
}

func (repo *serviceRepository) Save(ctx context.Context, service *domain.Service) (*domain.Service, error) {
	serviceDB, err := repo.db.CreateService(ctx, newCreateServiceParams(service))
	if err != nil {
		return nil, err
	}

	serviceDomain, err := newService(serviceDB)
	if err != nil {
		return nil, err
	}

	return serviceDomain, nil
}

func newCreateServiceParams(service *domain.Service) db.CreateServiceParams {
	float := decimal.NewFromFloat(service.Price)
	price := pgtype.Numeric{
		Int:   float.BigInt(),
		Exp:   float.Exponent(),
		NaN:   false,
		Valid: true,
	}
	return db.CreateServiceParams{
		ID:       service.ID,
		TenantID: service.TenantID,
		Name:     service.Name,
		Price:    price,
		Description: pgtype.Text{
			String: service.Description,
			Valid:  true,
		},
	}
}

func newService(serviceDB db.Service) (*domain.Service, error) {
	d, err := pgNumericToDecimal(serviceDB.Price)
	if err != nil {
		return nil, err
	}
	price := decimalToFloat64(d)
	return &domain.Service{
		ID:          serviceDB.ID,
		TenantID:    serviceDB.TenantID,
		Name:        serviceDB.Name,
		Price:       price,
		Description: serviceDB.Description.String,
	}, nil
}
