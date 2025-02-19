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
)

var interruptSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if cfg.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	connPool, err := pgxpool.New(ctx, cfg.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to database")
	}
	defer connPool.Close()

	runMigrations(cfg.MigrationURL, cfg.DBSource)

	store := db.NewStore(connPool)

	server, err := api.NewServer(store, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	srv := &http.Server{
		Addr:    cfg.HTTPServerAddr,
		Handler: server.Handler(),
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("cannot start server")
		}
	}()
	<-ctx.Done()
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
