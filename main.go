package main

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jpmoraess/pay/api"
	"github.com/jpmoraess/pay/config"
	db "github.com/jpmoraess/pay/db/sqlc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var interruptSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	// setup logger
	setupLogger(cfg.Environment)

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	// database
	connPool := mustInitDatabase(ctx, cfg)
	defer func() {
		log.Info().Msg("closing database connection")
		connPool.Close()
	}()

	runMigrations(cfg.MigrationURL, cfg.DBSource)

	store := db.NewStore(connPool)

	// create HTTP server
	srv := mustInitServer(store, cfg)

	errCh := make(chan error, 1)

	go func() {
		log.Info().Msgf("starting server on %s", srv.Addr)
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// graceful shutdown
	gracefulShutdown(ctx, srv, errCh)
}

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

func setupLogger(env string) {
	if env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

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

func mustInitServer(store db.Store, cfg *config.Config) *http.Server {
	server, err := api.NewServer(store, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}
	return &http.Server{
		Addr:    cfg.HTTPServerAddr,
		Handler: server.Handler(),
	}
}

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
