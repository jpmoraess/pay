package http

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jpmoraess/pay/config"
	"github.com/jpmoraess/pay/internal/application/ports"
	"github.com/jpmoraess/pay/internal/application/usecases"
	"github.com/jpmoraess/pay/token"
	"net/http"
	"time"
)

type TokenHandler struct {
	cfg            *config.Config
	tokenMaker     token.Maker
	userService    ports.UserService
	sessionService ports.SessionService
}

func NewTokenHandler(cfg *config.Config, tokenMaker token.Maker, router *gin.Engine, userService ports.UserService, sessionService ports.SessionService) {
	handler := TokenHandler{cfg, tokenMaker, userService, sessionService}
	group := router.Group("/tokens")
	{
		group.POST("/renew", handler.RenewToken)
	}
}

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

// RenewToken Renew access token
//
//	@Summary		Renew access token
//	@Description	Perform access token renew
//	@Tags			tokens
//	@Accept			json
//	@Produce		json
//	@Param			request	body		renewAccessTokenRequest	true	"Token renew request data"
//	@Success		200		{object}	renewAccessTokenResponse
//	@Failure		400		{object}	map[string]string
//	@Router			/tokens/renew [post]
func (h *TokenHandler) RenewToken(c *gin.Context) {
	var req renewAccessTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshPayload, err := h.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	session, err := h.sessionService.GetSession(c.Request.Context(), refreshPayload.ID)
	if err != nil {
		if errors.Is(err, usecases.ErrSessionNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if session.IsBlocked {
		err = fmt.Errorf("blocked session")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if session.Email != refreshPayload.Email {
		err = fmt.Errorf("incorrect session user")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err = fmt.Errorf("mismatched session token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err = fmt.Errorf("expired session")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	accessToken, accessTokenPayload, err := h.tokenMaker.CreateToken(
		refreshPayload.TenantID,
		refreshPayload.Email,
		refreshPayload.Role,
		h.cfg.AccessTokenDuration,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenPayload.ExpiredAt,
	}

	c.JSON(http.StatusOK, response)
}
