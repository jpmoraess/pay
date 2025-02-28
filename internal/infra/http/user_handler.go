package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jpmoraess/pay/internal/application/ports"
	"github.com/jpmoraess/pay/internal/application/usecases"
	"github.com/jpmoraess/pay/internal/infra/middleware"
	"github.com/jpmoraess/pay/token"
	"net/http"
	"time"
)

type UserHandler struct {
	tokenMaker  token.Maker
	userService ports.UserService
}

func NewUserHandler(router *gin.Engine, tokenMaker token.Maker, userService ports.UserService) {
	handler := &UserHandler{userService: userService}
	group := router.Group("/users").
		Use(middleware.AuthMiddleware(tokenMaker))
	{
		group.POST("/", handler.CreateUser).Use(middleware.RoleRequired([]string{"admin"}))
	}
}

type createUserRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@doe.com"`
	Password string `json:"password" binding:"required,min=8,max=20" example:"123456"`
	FullName string `json:"full_name" binding:"required,max=50" example:"John Doe"`
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateUser Create user
//
//	@Summary		Create user
//	@Description	Perform user creation
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		createUserRequest	true	"create request data"
//	@Success		201		{object}	userResponse
//	@Failure		400		{object}	map[string]string
//	@Router			/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req createUserRequest
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

	output, err := h.userService.Create(c.Request.Context(), &ports.CreateUserInput{
		TenantID: auth.TenantID,
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
		Role:     "service_provider",
	})
	if err != nil {
		if errors.Is(err, usecases.ErrUserAlreadyExists) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := userResponse{
		ID:        output.ID,
		Email:     output.Email,
		FullName:  output.FullName,
		CreatedAt: output.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}
