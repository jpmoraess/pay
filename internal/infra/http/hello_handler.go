package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jpmoraess/pay/internal/infra/middleware"
	"github.com/jpmoraess/pay/token"
	"net/http"
)

type HelloHandler struct {
	tokenMaker token.Maker
}

func NewHelloHandler(router *gin.Engine, tokenMaker token.Maker) {
	handler := &HelloHandler{}
	helloGroup := router.Group("/hello").
		Use(middleware.AuthMiddleware(tokenMaker))

	helloGroup.GET("/", handler.HelloWorld)
}

func (h *HelloHandler) HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "hello world"})
}
