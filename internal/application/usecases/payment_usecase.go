package usecases

import (
	"context"
	"github.com/jpmoraess/pay/internal/application/ports"
	"github.com/jpmoraess/pay/internal/domain"
	"github.com/jpmoraess/pay/internal/infra/gateway"
)

type paymentUseCase struct {
	asaas      *gateway.Asaas // TODO: use payment interface instead of concrete class
	repository ports.PaymentRepository
}

func NewPaymentUseCase(repository ports.PaymentRepository, assas *gateway.Asaas) *paymentUseCase {
	return &paymentUseCase{repository: repository, asaas: assas}
}

func (uc *paymentUseCase) Create(ctx context.Context, input *ports.CreatePaymentInput) (*ports.CreatePaymentOutput, error) {
	payment, err := domain.NewPayment(input.Value, input.DueDate)
	if err != nil {
		return nil, err
	}

	paymentResponse, err := uc.asaas.CreatePayment(ctx, &gateway.CreatePaymentRequest{
		Customer:          "6347643",
		BillingType:       gateway.Pix,
		Value:             payment.Value,
		DueDate:           payment.DueDate.Format("2006-01-02"),
		Description:       "payment generated go pay",
		ExternalReference: payment.ID.String(),
	})
	if err != nil {
		return nil, err
	}

	// set payment external id reference
	payment.ExternalID = paymentResponse.ID

	payment, err = uc.repository.Save(ctx, payment)
	if err != nil {
		return nil, err
	}

	output := &ports.CreatePaymentOutput{
		ID:    payment.ID,
		Value: payment.Value,
	}

	return output, nil
}
