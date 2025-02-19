package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jpmoraess/pay/token"
	"net/http"
	"strings"
)

const (
	authType       = "Bearer"
	authHeaderKey  = "Authorization"
	authPayloadKey = "auth_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authHeaderKey)
		if len(authHeader) == 0 {
			err := errors.New("authorization header is not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authType {
			err := fmt.Errorf("unsupported authorization type: %s", authorizationType)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		c.Set(authPayloadKey, payload)
		c.Next()
	}
}
