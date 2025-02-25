package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jpmoraess/pay/internal/application/ports"
	"github.com/jpmoraess/pay/internal/application/usecases"
	"net/http"
)

type TenantHandler struct {
	tenantService ports.TenantService
}

func NewTenantHandler(router *gin.Engine, tenantService ports.TenantService) {
	handler := &TenantHandler{tenantService: tenantService}
	group := router.Group("/tenants")
	{
		group.POST("/", handler.RegisterTenant)
	}
}

type registerTenantRequest struct {
	Email    string `json:"email" binding:"required" example:"john_doe@mail.com"`
	Name     string `json:"name" binding:"required" example:"John Doe Barber"`
	FullName string `json:"full_name" binding:"required" example:"John Doe Silva"`
	Password string `json:"password" binding:"required,min=8,max=20" example:"secretPwd"`
}

type registerTenantResponse struct {
	ID uuid.UUID `json:"id"`
}

// RegisterTenant Register tenant
//
//	@Summary		Register tenant
//	@Description	Perform tenant registration
//	@Tags			tenants
//	@Accept			json
//	@Produce		json
//	@Param			request	body		registerTenantRequest	true	"register request data"
//	@Success		201		{object}	registerTenantResponse
//	@Failure		400		{object}	map[string]string
//	@Router			/tenants [post]
func (h *TenantHandler) RegisterTenant(c *gin.Context) {
	var req registerTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.tenantService.Register(c.Request.Context(), &ports.RegisterTenantInput{
		Name:     req.Name,
		Email:    req.Email,
		FullName: req.FullName,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, usecases.ErrTenantAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := registerTenantResponse{ID: output.ID}

	c.JSON(http.StatusCreated, response)
}
