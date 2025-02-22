package domain

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type paymentStatus int

const (
	Pending paymentStatus = iota
	Paid
	Failed
	Cancelled
)

var paymentStatusToString = map[paymentStatus]string{
	Pending:   "PENDING",
	Paid:      "PAID",
	Failed:    "FAILED",
	Cancelled: "CANCELLED",
}

type Payment struct {
	ID         uuid.UUID
	ExternalID string
	Value      float64
	DueDate    time.Time
	Status     paymentStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewPayment(value float64, dueDate time.Time) (*Payment, error) {
	payment := &Payment{
		ID:         uuid.New(),
		ExternalID: uuid.New().String(),
		Value:      value,
		DueDate:    dueDate,
		Status:     Pending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := payment.validate(); err != nil {
		return nil, err
	}

	return payment, nil
}

func (p *Payment) validate() error {
	if p.ID == uuid.Nil {
		return errors.New("payment ID is required")
	}

	if len(p.ExternalID) == 0 {
		return errors.New("payment external ID is required")
	}

	if p.DueDate.Before(time.Now()) {
		return errors.New("payment due date must be in the future")
	}

	if len(p.ExternalID) == 0 {
		return errors.New("payment external ID is required")
	}

	if p.Value == 0 {
		return errors.New("payment value must not be zero")
	}
	return nil
}
