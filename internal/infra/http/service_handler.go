package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jpmoraess/pay/internal/application/ports"
	"github.com/jpmoraess/pay/internal/infra/middleware"
	"github.com/jpmoraess/pay/token"
	"net/http"
)

type ServiceHandler struct {
	tokenMaker     token.Maker
	serviceService ports.ServiceService
}

func NewServiceHandler(router *gin.Engine, tokenMaker token.Maker, serviceService ports.ServiceService) {
	handler := &ServiceHandler{tokenMaker: tokenMaker, serviceService: serviceService}
	group := router.Group("/services").
		Use(middleware.AuthMiddleware(tokenMaker))
	{
		group.POST("/", handler.CreateService)
	}
}

type createServiceRequest struct {
	Name        string  `json:"name" binding:"required" example:"Haircut"`
	Price       float64 `json:"price" binding:"required,gt=0" example:"29.99"`
	Description string  `json:"description" binding:"required" example:"This is a description"`
}

type createServiceResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
}

// CreateService Create service
//
//	@Summary		Create service
//	@Description	Perform service creation
//	@Tags			services
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		createServiceRequest	true	"Service request data"
//	@Success		201		{object}	createServiceResponse
//	@Failure		400		{object}	map[string]string
//	@Router			/services [post]
func (h *ServiceHandler) CreateService(c *gin.Context) {
	var req createServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	auth, ok := c.MustGet("auth_payload").(*token.Payload)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		c.Abort()
		return
	}

	output, err := h.serviceService.Create(c.Request.Context(), &ports.CreateServiceInput{
		TenantID:    auth.TenantID,
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := createServiceResponse{
		ID:          output.ID,
		Name:        output.Name,
		Price:       output.Price,
		Description: output.Description,
	}

	c.JSON(http.StatusCreated, response)
}
