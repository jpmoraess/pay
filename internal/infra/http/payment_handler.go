package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jpmoraess/pay/internal/application/ports"
	"github.com/jpmoraess/pay/internal/infra/middleware"
	"github.com/jpmoraess/pay/token"
	"net/http"
	"time"
)

type PaymentHandler struct {
	tokenMaker     token.Maker
	paymentService ports.PaymentService
}

func NewPaymentHandler(router *gin.Engine, tokenMaker token.Maker, paymentService ports.PaymentService) {
	handler := &PaymentHandler{tokenMaker: tokenMaker, paymentService: paymentService}
	group := router.Group("/payments").
		Use(middleware.AuthMiddleware(tokenMaker))
	{
		group.POST("/", handler.CreatePayment)
	}
}

type createPaymentRequest struct {
	Value   float64 `json:"value" binding:"required,gt=0" example:"20.99"`
	DueDate string  `json:"due_date" binding:"required" example:"2025-02-26"`
}

type createPaymentResponse struct {
	ID    uuid.UUID `json:"id"`
	Value float64   `json:"value"`
}

// CreatePayment Create payment
//
//	@Summary		Create payment
//	@Description	Perform payment creation
//	@Tags			payments
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		createPaymentRequest	true	"Payment request data"
//	@Success		201		{object}	createPaymentResponse
//	@Failure		400		{object}	map[string]string
//	@Router			/payments [post]
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req createPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parsedDueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.paymentService.Create(c, &ports.CreatePaymentInput{
		Value:   req.Value,
		DueDate: parsedDueDate,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := createPaymentResponse{
		ID:    output.ID,
		Value: req.Value,
	}

	c.JSON(http.StatusCreated, response)
}
