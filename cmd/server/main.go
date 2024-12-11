package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/northmule/gophkeeper/db"
	"github.com/northmule/gophkeeper/internal/server/api/handlers"
	"github.com/northmule/gophkeeper/internal/server/config"
	"github.com/northmule/gophkeeper/internal/server/logger"
	"github.com/northmule/gophkeeper/internal/server/repository"
	"github.com/northmule/gophkeeper/internal/server/storage"
)

func main() {
	fmt.Println("Running server gophkeeper...")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	var err error

	cfg := config.NewConfig()
	err = cfg.Init()
	if err != nil {
		return err
	}
	log, err := logger.NewLogger(cfg.Value().LogLevel)
	if err != nil {
		return err
	}
	log.Infof("Application Configuration: %#v", cfg.Value())

	log.Info("Database initialization")
	store, err := storage.NewPostgres(cfg.Value().Dsn)
	if err != nil {
		return err
	}

	log.Info("Checking the connection to the database")
	err = store.Ping(ctx)
	if err != nil {
		return err
	}

	log.Info("Initializing migrations")
	migrations := db.NewMigrations(store.RawDB)
	err = migrations.Up(ctx)
	if err != nil {
		return err
	}

	log.Info("Initializing the Repository Manager")
	repositoryManager, err := repository.NewManager(store.DB)
	if err != nil {
		return err
	}

	log.Info("Preparing the server for launch")
	routes := handlers.NewAppRoutes(repositoryManager, store.DB, storage.NewSession(), log, cfg)
	httpServer := http.Server{
		Addr:    cfg.Value().Address,
		Handler: routes.DefiningAppRoutes(),
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}
	go func() {
		<-ctx.Done()
		log.Info("The signal was received. I'm stopping the server...")

		shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()
		err = httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Error(err)
		}
	}()

	log.Infof("Running server on - %s", cfg.Value().Address)
	err = httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	if errors.Is(err, http.ErrServerClosed) {
		log.Info("The server is stopped")
	}

	return nil
}
