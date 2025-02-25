package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jpmoraess/pay/config"
	db "github.com/jpmoraess/pay/db/sqlc"
	"github.com/jpmoraess/pay/docs"
	_ "github.com/jpmoraess/pay/docs"
	"github.com/jpmoraess/pay/internal/adapters/database"
	"github.com/jpmoraess/pay/internal/application/ports"
	"github.com/jpmoraess/pay/internal/application/usecases"
	"github.com/jpmoraess/pay/internal/infra/gateway"
	handlers "github.com/jpmoraess/pay/internal/infra/http"
	"github.com/jpmoraess/pay/token"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var interruptSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}

//	@title			Pay
//	@version		1.0.0
//	@description	PayGolang

// @securityDefinitions.apiKey	BearerAuth
// @type						apiKey
// @in							header
// @name						Authorization
// @description				Enter "Bearer {your_token}" in the field below (without quotes)
func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	swaggerHost := cfg.SwaggerHost
	if len(swaggerHost) == 0 {
		swaggerHost = "localhost:8080"
	}
	docs.SwaggerInfo.Host = swaggerHost

	setupLogger(cfg.Environment)

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	connPool := mustInitDatabase(ctx, cfg)
	defer func() {
		log.Info().Msg("closing database connection")
		connPool.Close()
	}()

	runMigrations(cfg.MigrationURL, cfg.DBSource)

	store := db.NewStore(connPool)

	tokenMaker, err := token.NewPasetoMaker(cfg.SymmetricKey)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create token maker")
	}

	asaas := gateway.NewAsaas(cfg, &http.Client{
		Timeout: time.Second * 10,
	})

	// tenants
	tenantRepository := database.NewTenantRepository(store)
	tenantUseCase := usecases.NewTenantUseCase(tenantRepository)

	// sessions
	sessionRepository := database.NewSessionRepository(store)
	sessionUseCase := usecases.NewSessionUseCase(sessionRepository)

	// users
	userRepository := database.NewUserRepository(store)
	userUseCase := usecases.NewUserUseCase(cfg, tokenMaker, userRepository, sessionUseCase)

	// payments
	paymentRepository := database.NewPaymentRepository(store)
	paymentUseCase := usecases.NewPaymentUseCase(paymentRepository, asaas)

	// services
	serviceRepository := database.NewServiceRepository(store)
	serviceUseCase := usecases.NewServiceUseCase(serviceRepository)

	router := setupRouter(cfg, tokenMaker, userUseCase, tenantUseCase, sessionUseCase, paymentUseCase, serviceUseCase)

	srv := mustInitServer(router, cfg)

	errCh := make(chan error, 1)

	go func() {
		log.Info().Msgf("starting server on %s", srv.Addr)
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	gracefulShutdown(ctx, srv, errCh)
}

// runMigrations - run database migrations
func runMigrations(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("database migrated successfully")
}

// setupLogger - setup logger according environment
func setupLogger(env string) {
	if env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

// mustInitDatabase - init database connection
func mustInitDatabase(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	poolConfig, err := pgxpool.ParseConfig(cfg.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid database configuration")
	}

	connPool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to database")
	}

	return connPool
}

// mustInitServer - initiate http server
func mustInitServer(router *gin.Engine, cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:    cfg.HTTPServerAddr,
		Handler: router,
	}
}

// gracefulShutdown - shutdown http server gracefully
func gracefulShutdown(ctx context.Context, src *http.Server, errCh <-chan error) {
	select {
	case <-ctx.Done():
		log.Info().Msg("shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := src.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("server forced to shutdown")
		} else {
			log.Info().Msg("server gracefully shutdown")
		}

	case err := <-errCh:
		log.Error().Err(err).Msg("server encountered an error")
	}
}

// setupRouter - init gin and configure middlewares and handlers
func setupRouter(
	cfg *config.Config,
	tokenMaker token.Maker,
	userService ports.UserService,
	tenantUseCase ports.TenantService,
	sessionService ports.SessionService,
	paymentService ports.PaymentService,
	serviceUseCase ports.ServiceService,
) *gin.Engine {
	router := gin.Default()
	setupMiddlewares(router, cfg)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// asaas webhook route
	gateway.NewAsaasWebhook(router, cfg)

	handlers.NewHelloHandler(router, tokenMaker)
	handlers.NewUserHandler(router, userService)
	handlers.NewTenantHandler(router, tenantUseCase)
	handlers.NewPaymentHandler(router, tokenMaker, paymentService)
	handlers.NewServiceHandler(router, tokenMaker, serviceUseCase)
	handlers.NewTokenHandler(cfg, tokenMaker, router, userService, sessionService)

	return router
}

// setupMiddlewares - configure the gin middlewares
func setupMiddlewares(router *gin.Engine, cfg *config.Config) {
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.Info().
			Str("status", fmt.Sprintf("%d", param.StatusCode)).
			Str("method", param.Method).
			Str("path", param.Path).
			Dur("latency", param.Latency).
			Str("client_ip", param.ClientIP).
			Msg("incoming request")
		return ""
	}))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	router.Use(gzip.Gzip(gzip.DefaultCompression))
}
