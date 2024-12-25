package view

import (
	"testing"

	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestNewClientView(t *testing.T) {
	cfg := config.NewConfig()
	log, _ := logger.NewLogger("info")

	clientView := NewClientView(cfg, log)
	assert.NotNil(t, clientView)
	assert.Equal(t, cfg, clientView.cfg)
	assert.Equal(t, log, clientView.log)
}

func TestInitMain(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = "testpath"
	mockCfg.Value().PathKeys = "testpath"
	log, _ := logger.NewLogger("info")

	clientView := NewClientView(mockCfg, log)
	err := clientView.InitMain(context.Background())
	assert.EqualError(t, err, "open testpath/public_key.pem: no such file or directory")
}
