package ports

import (
	"context"
	"github.com/google/uuid"
	"github.com/jpmoraess/pay/internal/domain"
	"time"
)

type CreatePaymentInput struct {
	Value   float64
	DueDate time.Time
}

type CreatePaymentOutput struct {
	ID    uuid.UUID
	Value float64
}

type PaymentService interface {
	Create(ctx context.Context, input *CreatePaymentInput) (*CreatePaymentOutput, error)
}

type PaymentRepository interface {
	Save(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
}
