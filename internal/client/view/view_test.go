package view

import (
	"crypto/rand"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/common/keys"
	"github.com/northmule/gophkeeper/internal/common/keys/signers"
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

func TestInitMain_error(t *testing.T) {
	log, _ := logger.NewLogger("info")
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")

	os.MkdirAll(mockCfg.Value().PathKeys, 0755)
	os.WriteFile(filepath.Join(mockCfg.Value().PathKeys, keys.PrivateKeyFileNameForEncryption), []byte("encryption_key"), 0644)

	serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	kk := keys.NewKeys(keys.Options{
		Generator:    signers.NewRsaSigner(),
		SavePath:     mockCfg.Value().PathKeys,
		Organization: "Go32_client",
		Country:      "RU",
		SerialNumber: serialNumber,
	})
	_ = kk.InitSelfSigned()

	defer os.RemoveAll("testpath")

	clientView := NewClientView(mockCfg, log)
	ctx := context.Background()
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()
	err := clientView.InitMain(ctx)
	assert.NotEmpty(t, err)
}
