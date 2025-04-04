package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jpmoraess/pay/config"
	"io"
	"net/http"
	"strings"
	"time"
)

type Asaas struct {
	cfg    *config.Config
	client *http.Client
}

type BillingType int

const (
	Pix BillingType = iota
	Boleto
	CreditCard
)

var BillingTypeToString = map[BillingType]string{
	Pix:        "PIX",
	Boleto:     "BOLETO",
	CreditCard: "CREDIT_CARD",
}

type CreatePaymentRequest struct {
	Customer          string      `json:"customer"`
	BillingType       BillingType `json:"billingType"`
	Value             float64     `json:"value"`
	DueDate           string      `json:"dueDate"`
	Description       string      `json:"description"`
	ExternalReference string      `json:"externalReference"`
}

type CreatePaymentResponse struct {
	ID        string  `json:"id"`
	CreatedAt string  `json:"dateCreated"`
	Value     float64 `json:"value"`
}

type ErrorResponse struct {
	Errors []struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"errors"`
}

func NewAsaas(cfg *config.Config, client *http.Client) *Asaas {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &Asaas{cfg, client}
}

func (a *Asaas) CreatePayment(ctx context.Context, request *CreatePaymentRequest) (response *CreatePaymentResponse, err error) {
	url := fmt.Sprintf("%s/v3/payments", strings.TrimRight(a.cfg.AsaasURL, "/"))

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create payment request: %w", err)
	}

	b := bytes.NewReader(body)

	req, err := http.NewRequestWithContext(ctx, "POST", url, b)
	if err != nil {
		return nil, fmt.Errorf("failed to create create payment request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "jpmoraess/pay")
	req.Header.Set("access_token", a.cfg.AsaasApiKey)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute create payment request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read create payment response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorResponse ErrorResponse
		if err = json.Unmarshal(respBody, &errorResponse); err != nil {
			return nil, fmt.Errorf("error parsing error response (status: %d): %s", resp.StatusCode, string(respBody))
		}

		var errorMessages []string
		for _, e := range errorResponse.Errors {
			errorMessages = append(errorMessages, fmt.Sprintf("code: %s, description: %s", e.Code, e.Description))
		}
		return nil, fmt.Errorf("Asaas API error: %s", errorMessages)
	}

	if err = json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create payment response body: %w", err)
	}

	return response, nil
}

func (b BillingType) MarshalJSON() ([]byte, error) {
	if str, ok := BillingTypeToString[b]; ok {
		return json.Marshal(str)
	}
	return nil, fmt.Errorf("invalid billingType: %d", b)
}

func (b *BillingType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	for key, val := range BillingTypeToString {
		if val == s {
			*b = key
			return nil
		}
	}
	return fmt.Errorf("invalid billingType: %s", s)
}
