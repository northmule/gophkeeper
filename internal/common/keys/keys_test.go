package keys

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"testing"
)

type MockKeyGenerator struct {
	Key crypto.Signer
	Err error
}

func (m *MockKeyGenerator) GenerateKey() (crypto.Signer, error) {
	return m.Key, m.Err
}

func TestInitSelfSigned(t *testing.T) {
	testDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	mockGenerator := &MockKeyGenerator{Key: privateKey}
	options := Options{
		Generator:    mockGenerator,
		SavePath:     testDir,
		Organization: "TestOrg",
		Country:      "TestCountry",
		SerialNumber: big.NewInt(1),
	}

	keysService := NewKeys(options)

	err = keysService.InitSelfSigned()
	if err != nil {
		t.Errorf("InitSelfSigned failed: %v", err)
	}

	privateKeyPath := filepath.Join(testDir, PrivateKeyFileName)
	publicKeyPath := filepath.Join(testDir, PublicKeyFileName)
	certPath := filepath.Join(testDir, CertificateFileName)

	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		t.Errorf("Private key file not created: %v", err)
	}

	if _, err := os.Stat(publicKeyPath); os.IsNotExist(err) {
		t.Errorf("Public key file not created: %v", err)
	}

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		t.Errorf("Certificate file not created: %v", err)
	}

	mockGenerator.Err = errors.New("mock error")
	err = keysService.InitSelfSigned()
	if err == nil {
		t.Errorf("InitSelfSigned should have failed with mock error")
	}

	mockGenerator.Err = nil
	mockGenerator.Key = nil
	err = keysService.InitSelfSigned()
	if err == nil {
		t.Errorf("InitSelfSigned should have failed with nil key")
	}

	// Test input validation
	options.SavePath = ""
	keysService = NewKeys(options)
	if keysService.PrivateKeyPath() != PrivateKeyFileName {
		t.Errorf("PrivateKeyPath should be empty with empty SavePath")
	}
	if keysService.PublicKeyPath() != PublicKeyFileName {
		t.Errorf("PublicKeyFileName should be empty with empty SavePath")
	}
	if keysService.CertPath() != CertificateFileName {
		t.Errorf("CertificateFileName should be empty with empty SavePath")
	}
}
