package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/common/keys"
	"github.com/northmule/gophkeeper/internal/common/keys/signers"
	"github.com/stretchr/testify/assert"
)

func generateRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey
}

func writeKeyToFile(key interface{}, filename string) {
	var pemBlock *pem.Block
	switch k := key.(type) {
	case *rsa.PrivateKey:
		pemBlock = &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(k),
		}
	case *rsa.PublicKey:
		pemBlock = &pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(k),
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = pem.Encode(file, pemBlock)
	if err != nil {
		panic(err)
	}
}

func TestNewCrypt_Success(t *testing.T) {
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
	crypt, err := NewCrypt(mockCfg)
	assert.NoError(t, err)
	assert.NotNil(t, crypt)
	assert.NotNil(t, crypt.serverPublicKey)
	assert.NotNil(t, crypt.clientPrivateKey)
	assert.NotNil(t, crypt.privateKeyForEncryption)
	assert.Equal(t, "encryption_key", string(crypt.privateKeyForEncryption))

}

func TestNewCrypt_ServerPublicKeyFileError(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")

	os.MkdirAll(mockCfg.Value().PathKeys, 0755)

	crypt, err := NewCrypt(mockCfg)
	assert.Error(t, err)
	assert.Nil(t, crypt)

	os.RemoveAll(mockCfg.Value().PathKeys)
}

func TestNewCrypt_ClientPrivateKeyFileError(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")

	_, publicKey := generateRSAKeyPair()

	os.MkdirAll(mockCfg.Value().PathKeys, 0755)
	writeKeyToFile(publicKey, filepath.Join(mockCfg.Value().PathKeys, keys.PublicKeyFileName))

	crypt, err := NewCrypt(mockCfg)
	assert.Error(t, err)
	assert.Nil(t, crypt)

	os.RemoveAll(mockCfg.Value().PathKeys)
}

func TestNewCrypt_EncryptionKeyFileError(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")

	privateKey, publicKey := generateRSAKeyPair()

	os.MkdirAll(mockCfg.Value().PathKeys, 0755)
	writeKeyToFile(publicKey, filepath.Join(mockCfg.Value().PathKeys, keys.PublicKeyFileName))
	writeKeyToFile(privateKey, filepath.Join(mockCfg.Value().PathKeys, keys.PrivateKeyFileName))

	crypt, err := NewCrypt(mockCfg)
	assert.Error(t, err)
	assert.Nil(t, crypt)

	os.RemoveAll(mockCfg.Value().PathKeys)
}

func TestCrypt_EncryptRSA(t *testing.T) {
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
	crypt, _ := NewCrypt(mockCfg)
	v, e := crypt.EncryptRSA([]byte("text"))
	assert.NoError(t, e)
	assert.NotEmpty(t, v)
}
func TestCrypt_DecryptRSA(t *testing.T) {
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
	crypt, _ := NewCrypt(mockCfg)
	vv, _ := crypt.EncryptRSA([]byte("text"))
	v, e := crypt.DecryptRSA(vv)
	assert.NoError(t, e)
	assert.NotEmpty(t, v)
}

func TestCrypt_EncryptAES(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	key := make([]byte, 32) // AES-256 key
	os.MkdirAll(mockCfg.Value().PathKeys, 0755)
	os.WriteFile(filepath.Join(mockCfg.Value().PathKeys, keys.PrivateKeyFileNameForEncryption), key, 0644)

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
	v, e := crypt.EncryptAES([]byte("text"))
	assert.NoError(t, e)
	assert.NotEmpty(t, v)
}

func TestCrypt_DecryptAES(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")

	key := make([]byte, 32) // AES-256 key

	os.MkdirAll(mockCfg.Value().PathKeys, 0755)
	os.WriteFile(filepath.Join(mockCfg.Value().PathKeys, keys.PrivateKeyFileNameForEncryption), key, 0644)

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
	vv, _ := crypt.EncryptAES([]byte("text"))
	v, e := crypt.DecryptAES(vv)
	assert.NoError(t, e)
	assert.NotEmpty(t, v)
}
