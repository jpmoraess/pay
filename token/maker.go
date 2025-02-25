package token

import (
	"github.com/google/uuid"
	"time"
)

type Maker interface {
	CreateToken(tenantID uuid.UUID, email string, role string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
