package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jpmoraess/pay/config"
	db "github.com/jpmoraess/pay/db/sqlc"
	"github.com/jpmoraess/pay/token"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     *config.Config
}

func NewServer(store db.Store, config *config.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hello world"})
	})

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.GET("/hello-secure", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hello world secure"})
	})

	server.router = router
}

func (server *Server) Handler() *gin.Engine {
	return server.router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
