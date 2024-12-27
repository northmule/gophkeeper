package service

import (
	"crypto/rand"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/northmule/gophkeeper/internal/common/keys"
	"github.com/northmule/gophkeeper/internal/common/keys/signers"
	"github.com/northmule/gophkeeper/internal/server/config"
	"github.com/stretchr/testify/assert"
)

func TestNewCrypt_Success(t *testing.T) {
	mockCfg := config.NewConfig()
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
	crypt, err := NewCrypt(mockCfg)
	assert.NoError(t, err)
	assert.NotNil(t, crypt)
	assert.NotNil(t, crypt.serverPrivateKey)
	assert.NotNil(t, crypt.clientPublicKey)
}

func TestCrypt_EncryptRSA(t *testing.T) {
	mockCfg := config.NewConfig()
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
	crypt, _ := NewCrypt(mockCfg)
	v, e := crypt.EncryptRSA([]byte("text"))
	assert.NoError(t, e)
	assert.NotEmpty(t, v)
}
func TestCrypt_DecryptRSA(t *testing.T) {
	mockCfg := config.NewConfig()
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
	crypt, _ := NewCrypt(mockCfg)
	vv, _ := crypt.EncryptRSA([]byte("text"))
	v, e := crypt.DecryptRSA(vv)
	assert.NoError(t, e)
	assert.NotEmpty(t, v)
}
