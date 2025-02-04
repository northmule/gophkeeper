package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/northmule/gophkeeper/db"
	"github.com/northmule/gophkeeper/internal/common/keys"
	"github.com/northmule/gophkeeper/internal/common/keys/signers"
	"github.com/northmule/gophkeeper/internal/server/api/handlers"
	"github.com/northmule/gophkeeper/internal/server/config"
	"github.com/northmule/gophkeeper/internal/server/logger"
	"github.com/northmule/gophkeeper/internal/server/repository"
	service "github.com/northmule/gophkeeper/internal/server/services"
	"github.com/northmule/gophkeeper/internal/server/services/access"
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

	if cfg.Value().MigrationsApply {
		log.Info("Initializing migrations")
		migrations := db.NewMigrations(store.RawDB)
		err = migrations.Up(ctx)
		if err != nil {
			return err
		}
	} else {
		log.Info("Skip Initializing migrations")
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return err
	}

	log.Info("Preparing server keys")
	serverKeys := keys.NewKeys(keys.Options{
		Generator:    signers.NewRsaSigner(),
		SavePath:     cfg.Value().PathKeys,
		Organization: "Go32_Server",
		Country:      "RU",
		SerialNumber: serialNumber,
	})

	var errKeys []error
	if _, err = os.Stat(serverKeys.PrivateKeyPath()); errors.Is(err, os.ErrNotExist) {
		errKeys = append(errKeys, fmt.Errorf("private key file does not exist"))
	}
	if _, err = os.Stat(serverKeys.PublicKeyPath()); errors.Is(err, os.ErrNotExist) {
		errKeys = append(errKeys, fmt.Errorf("public key file does not exist"))
	}
	if _, err = os.Stat(serverKeys.CertPath()); errors.Is(err, os.ErrNotExist) {
		errKeys = append(errKeys, fmt.Errorf("cert file does not exist"))
	}

	if cfg.Value().OverwriteKeys || len(errKeys) > 0 {
		log.Info("Creating Server Keys")
		err = serverKeys.InitSelfSigned()
		if err != nil {
			return err
		}
	}

	log.Info("Preparing the server for launch")
	accessService := access.NewAccess(cfg)
	cryptService, err := service.NewCrypt(cfg)
	if err != nil {
		return err
	}

	log.Info("Initializing the Repository Manager")
	userRepository, err := repository.NewUserRepository(store.DB)
	if err != nil {
		return err
	}
	cardDataRepository, err := repository.NewCardDataRepository(store.DB)
	if err != nil {
		return err
	}
	ownerRepository, err := repository.NewOwnerRepository(store.DB)
	if err != nil {
		return err
	}
	metaDataRepository, err := repository.NewMetaDataRepository(store.DB)
	if err != nil {
		return err
	}
	textDataRepository, err := repository.NewTextDataRepository(store.DB)
	if err != nil {
		return err
	}
	fileDataRepository, err := repository.NewFileDataRepository(store.DB)
	if err != nil {
		return err
	}

	log.Info("Initializing the Routes")
	routes := handlers.NewAppRoutes(store.DB, storage.NewSession(), log, cfg, accessService, cryptService).
		SetFileDataRepository(fileDataRepository).
		SetCardDataRepository(cardDataRepository).
		SetMetaDataRepository(metaDataRepository).
		SetOwnerRepository(ownerRepository).
		SetTextDataRepository(textDataRepository).
		SetUserRepository(userRepository)

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
