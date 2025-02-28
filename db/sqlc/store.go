package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	CreateTenantTx(ctx context.Context, params CreateTenantTxParams) (CreateTenantTxResult, error)
}

type SQLStore struct {
	conPool *pgxpool.Pool
	*Queries
}

func NewStore(conPool *pgxpool.Pool) Store {
	return &SQLStore{
		conPool: conPool,
		Queries: New(conPool),
	}
}
