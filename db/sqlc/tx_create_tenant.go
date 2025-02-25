package db

import (
	"context"
	"github.com/google/uuid"
)

type CreateTenantTxParams struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type CreateTenantTxResult struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

func (store *SQLStore) CreateTenantTx(ctx context.Context, params CreateTenantTxParams) (CreateTenantTxResult, error) {
	var result CreateTenantTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		tenant, err := q.CreateTenant(ctx, CreateTenantParams{
			ID:       params.ID,
			Name:     params.Name,
			Email:    params.Email,
			Password: params.Password,
		})
		if err != nil {
			return err
		}

		_, err = q.CreateUser(ctx, CreateUserParams{
			ID:       tenant.ID,
			TenantID: tenant.ID,
			Email:    tenant.Email,
			FullName: params.FullName,
			Password: params.Password,
			Role:     "admin",
		})
		if err != nil {
			return err
		}

		result.ID = tenant.ID

		return err
	})

	return result, err
}
