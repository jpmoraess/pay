package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/jpmoraess/pay/db/sqlc"
	"github.com/jpmoraess/pay/internal/domain"
	"github.com/shopspring/decimal"
)

type paymentRepository struct {
	db db.Store
}

func NewPaymentRepository(db db.Store) *paymentRepository {
	return &paymentRepository{db: db}
}

func newPayment(paymentDB db.Payment) (*domain.Payment, error) {
	decimal, err := pgNumericToDecimal(paymentDB.Value)
	if err != nil {
		return nil, err
	}
	value := decimalToFloat64(decimal)
	return &domain.Payment{
		ID:         paymentDB.ID,
		ExternalID: paymentDB.ExternalID,
		Value:      value,
		DueDate:    paymentDB.DueDate,
		Status:     domain.Pending,
		CreatedAt:  paymentDB.CreatedAt,
		UpdatedAt:  paymentDB.UpdatedAt,
	}, nil
}

func newCreatePaymentParams(payment *domain.Payment) db.CreatePaymentParams {
	float := decimal.NewFromFloat(payment.Value)
	value := pgtype.Numeric{
		Int:   float.BigInt(),
		Exp:   float.Exponent(),
		NaN:   false,
		Valid: true,
	}
	return db.CreatePaymentParams{
		ID:         payment.ID,
		ExternalID: payment.ExternalID,
		Value:      value,
		DueDate:    payment.DueDate,
	}
}

func (repo *paymentRepository) Save(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	paymentDB, err := repo.db.CreatePayment(ctx, newCreatePaymentParams(payment))
	if err != nil {
		return nil, err
	}

	paymentDomain, err := newPayment(paymentDB)
	if err != nil {
		return nil, err
	}

	return paymentDomain, nil
}

// pgNumericToDecimal - convert pgtype to decimal
func pgNumericToDecimal(n pgtype.Numeric) (decimal.Decimal, error) {
	return decimal.NewFromBigInt(n.Int, n.Exp), nil
}

// decimalToFloat64 - convert decimal to float64
func decimalToFloat64(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}
