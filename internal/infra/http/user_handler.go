package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jpmoraess/pay/internal/application/ports"
	"github.com/jpmoraess/pay/internal/application/usecases"
	"net/http"
	"time"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(router *gin.Engine, userService ports.UserService) {
	handler := &UserHandler{userService: userService}
	group := router.Group("/users")
	{
		group.POST("/", handler.CreateUser)
		group.POST("/login", handler.UserLogin)
	}
}

type createUserRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@doe.com"`
	Password string `json:"password" binding:"required,min=8,max=20" example:"123456"`
	FullName string `json:"full_name" binding:"required,max=50" example:"John Doe"`
}

type userResponse struct {
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

	output, err := h.userService.Create(c.Request.Context(), &ports.CreateUserInput{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
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
		Email:     output.Email,
		FullName:  output.FullName,
		CreatedAt: output.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

type loginRequest struct {
	Email    string `json:"email" binding:"required" example:"john@doe.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type loginResponse struct {
	SessionID             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

// UserLogin User login
//
//	@Summary		User login
//	@Description	Perform user login
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		loginRequest	true	"Login request data"
//	@Success		200		{object}	loginResponse
//	@Failure		400		{object}	map[string]string
//	@Router			/users/login [post]
func (h *UserHandler) UserLogin(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Set("client-ip", c.ClientIP())
	c.Set("user-agent", c.Request.UserAgent())

	output, err := h.userService.Login(c.Request.Context(), &ports.LoginUserInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, usecases.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, usecases.ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := loginResponse{
		SessionID:             output.SessionID,
		AccessToken:           output.AccessToken,
		AccessTokenExpiresAt:  output.AccessTokenExpiresAt,
		RefreshToken:          output.RefreshToken,
		RefreshTokenExpiresAt: output.RefreshTokenExpiresAt,
	}

	c.JSON(http.StatusOK, response)
}
