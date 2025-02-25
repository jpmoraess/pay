package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jpmoraess/pay/token"
	"net/http"
)

func RoleRequired(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authPayload, ok := c.MustGet(authPayloadKey).(*token.Payload)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		roleValid := false
		for _, role := range allowedRoles {
			if role == authPayload.Role {
				roleValid = true
				break
			}
		}

		if !roleValid {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}
