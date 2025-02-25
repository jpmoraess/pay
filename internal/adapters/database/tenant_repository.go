package database

import (
	"context"
	db "github.com/jpmoraess/pay/db/sqlc"
	"github.com/jpmoraess/pay/internal/domain"
)

type tenantRepository struct {
	db db.Store
}

func NewTenantRepository(db db.Store) *tenantRepository {
	return &tenantRepository{db: db}
}

func (repo *tenantRepository) Save(ctx context.Context, tenant *domain.Tenant) (*domain.Tenant, error) {
	txResult, err := repo.db.CreateTenantTx(ctx, db.CreateTenantTxParams{
		ID:       tenant.ID,
		Name:     tenant.Name,
		FullName: tenant.FullName,
		Email:    tenant.Email,
		Password: tenant.Password,
	})
	if err != nil {
		return nil, err
	}
	return newTenant(txResult), nil
}

func newTenant(txResult db.CreateTenantTxResult) *domain.Tenant {
	return &domain.Tenant{
		ID:       txResult.ID,
		Name:     txResult.Name,
		FullName: txResult.FullName,
		Email:    txResult.Email,
		Password: txResult.Password,
	}
}
