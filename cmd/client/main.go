package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"context"

	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/controller"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/client/service"
	"github.com/northmule/gophkeeper/internal/client/storage"
	appview "github.com/northmule/gophkeeper/internal/client/view"
	"github.com/northmule/gophkeeper/internal/common/keys"
	"github.com/northmule/gophkeeper/internal/common/keys/signers"
)

var (
	version   = "develop"
	buildDate = "n/a"
)

func main() {
	fmt.Println("Running client gophkeeper...")

	fmt.Printf("Version: %s\n", version)
	fmt.Printf("BuildDate: %s\n", buildDate)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	var err error

	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	log, err := logger.NewLogger(cfg.Value().LogLevel)
	if err != nil {
		return err
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return err
	}
	clientKeys := keys.NewKeys(keys.Options{
		Generator:    signers.NewRsaSigner(),
		SavePath:     cfg.Value().PathKeys,
		Organization: "Go32_client",
		Country:      "RU",
		SerialNumber: serialNumber,
	})

	var errKeys []error
	if _, err = os.Stat(clientKeys.PrivateKeyPath()); errors.Is(err, os.ErrNotExist) {
		errKeys = append(errKeys, fmt.Errorf("private key file does not exist"))
	}
	if _, err = os.Stat(clientKeys.PublicKeyPath()); errors.Is(err, os.ErrNotExist) {
		errKeys = append(errKeys, fmt.Errorf("public key file does not exist"))
	}
	if _, err = os.Stat(clientKeys.CertPath()); errors.Is(err, os.ErrNotExist) {
		errKeys = append(errKeys, fmt.Errorf("cert file does not exist"))
	}

	if cfg.Value().OverwriteKeys || len(errKeys) > 0 {
		err = clientKeys.InitSelfSigned()
		if err != nil {
			return err
		}
	}

	cryptService, err := service.NewCrypt(cfg)
	if err != nil {
		return err
	}
	manager, err := controller.NewManager(cfg, cryptService, log)
	if err != nil {
		return err
	}

	clientView := appview.NewClientView(manager, storage.NewMemoryStorage(), log)

	return clientView.InitMain(ctx)
}
