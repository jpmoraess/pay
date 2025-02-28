// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: service.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createService = `-- name: CreateService :one
INSERT INTO services (id, tenant_id, name, price, description)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, tenant_id, name, price, description
`

type CreateServiceParams struct {
	ID          uuid.UUID      `json:"id"`
	TenantID    uuid.UUID      `json:"tenant_id"`
	Name        string         `json:"name"`
	Price       pgtype.Numeric `json:"price"`
	Description pgtype.Text    `json:"description"`
}

func (q *Queries) CreateService(ctx context.Context, arg CreateServiceParams) (Service, error) {
	row := q.db.QueryRow(ctx, createService,
		arg.ID,
		arg.TenantID,
		arg.Name,
		arg.Price,
		arg.Description,
	)
	var i Service
	err := row.Scan(
		&i.ID,
		&i.TenantID,
		&i.Name,
		&i.Price,
		&i.Description,
	)
	return i, err
}

const getService = `-- name: GetService :one
SELECT id, tenant_id, name, price, description FROM services
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetService(ctx context.Context, id uuid.UUID) (Service, error) {
	row := q.db.QueryRow(ctx, getService, id)
	var i Service
	err := row.Scan(
		&i.ID,
		&i.TenantID,
		&i.Name,
		&i.Price,
		&i.Description,
	)
	return i, err
}

const listServices = `-- name: ListServices :many
SELECT id, tenant_id, name, price, description FROM services
WHERE tenant_id = $1
`

func (q *Queries) ListServices(ctx context.Context, tenantID uuid.UUID) ([]Service, error) {
	rows, err := q.db.Query(ctx, listServices, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Service{}
	for rows.Next() {
		var i Service
		if err := rows.Scan(
			&i.ID,
			&i.TenantID,
			&i.Name,
			&i.Price,
			&i.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
