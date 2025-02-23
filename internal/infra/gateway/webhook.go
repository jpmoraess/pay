package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jpmoraess/pay/config"
	"net/http"
)

type PaymentEventType int

const (
	PaymentCreated PaymentEventType = iota
	PaymentConfirmed
	PaymentReceived
	PaymentOverdue
)

var PaymentEventTypeToString = map[PaymentEventType]string{
	PaymentCreated:   "PAYMENT_CREATED",
	PaymentConfirmed: "PAYMENT_CONFIRMED",
	PaymentReceived:  "PAYMENT_RECEIVED",
	PaymentOverdue:   "PAYMENT_OVERDUE",
}

type asaasPaymentEvent struct {
	ID          string `json:"id"`
	Event       string `json:"event"`
	DateCreated string `json:"dateCreated"`
	Payment     struct {
		Object          string  `json:"object"`
		ID              string  `json:"id"`
		DateCreated     string  `json:"dateCreated"`
		Customer        string  `json:"customer"`
		PaymentLink     string  `json:"paymentLink"`
		DueDate         string  `json:"dueDate"`
		OriginalDueDate string  `json:"originalDueDate"`
		Value           float64 `json:"value"`
		NetValue        float64 `json:"netValue"`
		BillingType     string  `json:"billingType"`
		Status          string  `json:"status"`
	} `json:"payment"`
}

type AsaasWebhook struct {
	config *config.Config
}

func NewAsaasWebhook(router *gin.Engine, config *config.Config) {
	webhook := &AsaasWebhook{config: config}
	group := router.Group("/asaas/webhooks")
	{
		group.POST("/payments", webhook.HandlePaymentEvent)
	}
}

func (wh *AsaasWebhook) HandlePaymentEvent(c *gin.Context) {
	// 1. receber o token do asaas
	token := c.GetHeader("asaas-access-token")
	_ = token

	// 2. validar access token com usuário do sistema

	// 3. realizar o parse do corpo da requisicao
	var req asaasPaymentEvent
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 4. processar evento
	switch req.Event {
	case "PAYMENT_CREATED":
		fmt.Printf("PAYMENT_CREATED: %s\n", req)
	case "PAYMENT_CONFIRMED":
		fmt.Printf("PAYMENT_CONFIRMED: %s\n", req)
	case "PAYMENT_OVERDUE":
		fmt.Printf("PAYMENT_OVERDUE: %s\n", req)
	case "PAYMENT_RECEIVED":
		fmt.Printf("PAYMENT_RECEIVED: %s\n", req)
	default:
		fmt.Println("UNKNOW EVENT TYPE")
	}

	// 5. retornar sucesso
	c.JSON(http.StatusOK, gin.H{"received": true})
}

func (p PaymentEventType) MarshalJSON() ([]byte, error) {
	if str, ok := PaymentEventTypeToString[p]; ok {
		return json.Marshal(str)
	}
	return nil, fmt.Errorf("invalid paymentEventType: %d", p)
}

func (p *PaymentEventType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	for key, val := range PaymentEventTypeToString {
		if val == s {
			*p = key
			return nil
		}
	}
	return fmt.Errorf("invalid paymentEventType: %s", s)
}
