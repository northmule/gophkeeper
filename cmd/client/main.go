package main

import (
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"context"

	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/logger"
	appview "github.com/northmule/gophkeeper/internal/client/view"
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

	cfg := config.NewConfig()
	err = cfg.Init()
	if err != nil {
		return err
	}
	log, err := logger.NewLogger(cfg.Value().LogLevel)
	if err != nil {
		return err
	}

	clientView := appview.NewClientView(log)

	return clientView.InitMain(ctx)
}
