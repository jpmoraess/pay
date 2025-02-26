package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jpmoraess/pay/internal/application/ports"
	"github.com/jpmoraess/pay/internal/application/usecases"
	"github.com/jpmoraess/pay/token"
	"net/http"
	"time"
)

type AuthHandler struct {
	tokenMaker  token.Maker
	userService ports.UserService
}

func NewAuthHandler(router *gin.Engine, tokenMaker token.Maker, userService ports.UserService) {
	handler := &AuthHandler{
		tokenMaker:  tokenMaker,
		userService: userService,
	}

	group := router.Group("/auth")
	{
		group.POST("/login", handler.AuthLogin)
	}
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

// AuthLogin Login
//
//	@Summary		Login
//	@Description	Perform login
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		loginRequest	true	"Login request data"
//	@Success		200		{object}	loginResponse
//	@Failure		400		{object}	map[string]string
//	@Router			/auth/login [post]
func (h *AuthHandler) AuthLogin(c *gin.Context) {
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
