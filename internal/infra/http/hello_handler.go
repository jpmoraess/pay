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
	handler := &HelloHandler{tokenMaker: tokenMaker}
	helloGroup := router.Group("/hello").
		Use(middleware.AuthMiddleware(tokenMaker))
	{
		helloGroup.GET("/", handler.HelloWorld)
	}
}

// HelloWorld Hello world
//
//	@Summary		Hello world
//	@Description	Perform hello world
//	@Tags			hello
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Router			/hello [get]
func (h *HelloHandler) HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "hello world"})
}
