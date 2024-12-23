package keys

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"testing"
)

// MockKeyGenerator is a mock implementation of KeyGenerator
type MockKeyGenerator struct {
	Key crypto.Signer
	Err error
}

func (m *MockKeyGenerator) GenerateKey() (crypto.Signer, error) {
	return m.Key, m.Err
}

func TestInitSelfSigned(t *testing.T) {

	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	mockGenerator := &MockKeyGenerator{Key: privateKey}
	options := Options{
		Generator:    mockGenerator,
		SavePath:     dir,
		Organization: "TestOrg",
		Country:      "TestCountry",
		SerialNumber: big.NewInt(1),
	}

	keysService := NewKeys(options)

	// Test main functionality
	err = keysService.InitSelfSigned()
	if err != nil {
		t.Errorf("InitSelfSigned failed: %v", err)
	}

	files := []string{
		PrivateKeyFileName,
		PublicKeyFileName,
		CertificateFileName,
	}

	for _, file := range files {
		path := filepath.Join(dir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("File %s does not exist", path)
		}
	}

	mockGenerator.Err = errors.New("mock error")
	err = keysService.InitSelfSigned()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	mockGenerator.Err = nil
	mockGenerator.Key = nil
	err = keysService.InitSelfSigned()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	options.SavePath = ""
	keysService = NewKeys(options)
	err = keysService.InitSelfSigned()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	options.SavePath = dir
	options.Organization = ""
	keysService = NewKeys(options)
	err = keysService.InitSelfSigned()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	options.Organization = "TestOrg"
	options.Country = ""
	keysService = NewKeys(options)
	err = keysService.InitSelfSigned()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	options.Country = "TestCountry"
	options.SerialNumber = big.NewInt(0)
	keysService = NewKeys(options)
	err = keysService.InitSelfSigned()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
