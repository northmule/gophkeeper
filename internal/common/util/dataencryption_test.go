package util

import (
	"crypto/aes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateHashForKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Normal String", "hello", "5d41402abc4b2a76b9719d911017c592"},
		{"Empty String", "", "d41d8cd98f00b204e9800998ecf8427e"},
		{"Whitespace String", "   ", "628631f07321b22d8c176c200c855e1b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateHashForKey(tt.input)
			if result != tt.expected {
				t.Errorf("CreateHashForKey(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFillPublicRsaKeyFromString(t *testing.T) {
	// Generate a valid RSA public key in PEM format
	rsaKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	rsaKeyBytes, _ := x509.MarshalPKIXPublicKey(&rsaKey.PublicKey)
	rsaKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: rsaKeyBytes,
	})

	t.Run("Valid RSA Key", func(t *testing.T) {
		key, err := FillPublicRsaKeyFromString(string(rsaKeyPem))
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if key == nil {
			t.Errorf("Expected non-nil key, got nil")
		}
	})

	t.Run("Empty String", func(t *testing.T) {
		key, err := FillPublicRsaKeyFromString("")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if key != nil {
			t.Errorf("Expected nil key, got non-nil")
		}
	})

	t.Run("Invalid PEM Data", func(t *testing.T) {
		key, err := FillPublicRsaKeyFromString("invalid pem data")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if key != nil {
			t.Errorf("Expected nil key, got non-nil")
		}
	})

	ecdsaKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	ecdsaKeyBytes, _ := x509.MarshalPKIXPublicKey(&ecdsaKey.PublicKey)
	ecdsaKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: ecdsaKeyBytes,
	})

	t.Run("Unsupported Key Type", func(t *testing.T) {
		key, err := FillPublicRsaKeyFromString(string(ecdsaKeyPem))
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if key != nil {
			t.Errorf("Expected nil key, got non-nil")
		}
	})
}

func TestFillPublicRsaKeyFromFile(t *testing.T) {

	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	validKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}
	validKeyBytes, err := x509.MarshalPKIXPublicKey(&validKey.PublicKey)
	if err != nil {
		t.Fatalf("Failed to marshal public key: %v", err)
	}
	validPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: validKeyBytes,
	})
	validFilePath := filepath.Join(tempDir, "valid.pem")
	if err := ioutil.WriteFile(validFilePath, validPem, 0644); err != nil {
		t.Fatalf("Failed to write valid PEM file: %v", err)
	}

	invalidPem := []byte("Invalid PEM data")
	invalidFilePath := filepath.Join(tempDir, "invalid.pem")
	if err := ioutil.WriteFile(invalidFilePath, invalidPem, 0644); err != nil {
		t.Fatalf("Failed to write invalid PEM file: %v", err)
	}

	emptyFilePath := filepath.Join(tempDir, "empty.pem")
	if _, err := os.Create(emptyFilePath); err != nil {
		t.Fatalf("Failed to create empty PEM file: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		expected error
	}{
		{"Valid RSA Key", validFilePath, nil},
		{"Invalid PEM Data", invalidFilePath, fmt.Errorf("no PEM data found in file")},
		{"Empty File", emptyFilePath, fmt.Errorf("no PEM data found in file")},
		{"Non-existent File", "nonexistent.pem", fmt.Errorf("open nonexistent.pem: no such file or directory")},
		{"Empty Path", "", fmt.Errorf("open : no such file or directory")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			key, err := FillPublicRsaKeyFromFile(test.path)
			if test.expected == nil {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if key == nil {
					t.Errorf("Expected non-nil key, got nil")
				}
			} else {
				if err == nil {
					t.Errorf("Expected error %v, got nil", test.expected)
				} else if err.Error() != test.expected.Error() {
					t.Errorf("Expected error %v, got %v", test.expected, err)
				}
			}
		})
	}
}

func createTempPemFile(t *testing.T, key *rsa.PrivateKey) string {
	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("Failed to marshal private key: %v", err)
	}

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	}

	file, err := ioutil.TempFile("", "testkey.pem")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	if err := pem.Encode(file, block); err != nil {
		t.Fatalf("Failed to encode PEM block: %v", err)
	}

	return file.Name()
}

func TestFillPrivateRsaKeyFromFile(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA private key: %v", err)
	}

	validPemPath := createTempPemFile(t, privateKey)

	invalidPemPath := filepath.Join(tempDir, "invalidkey.pem")
	if err := ioutil.WriteFile(invalidPemPath, []byte("invalid data"), 0644); err != nil {
		t.Fatalf("Failed to create invalid PEM file: %v", err)
	}

	emptyFilePath := filepath.Join(tempDir, "emptyfile.pem")
	if _, err := os.Create(emptyFilePath); err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	t.Run("ValidRSAKey", func(t *testing.T) {
		key, err := FillPrivateRsaKeyFromFile(validPemPath)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if key == nil {
			t.Errorf("Expected non-nil key, got nil")
		}
	})

	t.Run("InvalidRSAKey", func(t *testing.T) {
		_, err := FillPrivateRsaKeyFromFile(invalidPemPath)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("EmptyFile", func(t *testing.T) {
		_, err := FillPrivateRsaKeyFromFile(emptyFilePath)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("NonExistentFile", func(t *testing.T) {
		_, err := FillPrivateRsaKeyFromFile(filepath.Join(tempDir, "nonexistent.pem"))
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("EmptyFilePath", func(t *testing.T) {
		_, err := FillPrivateRsaKeyFromFile("")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("DirectoryPath", func(t *testing.T) {
		_, err := FillPrivateRsaKeyFromFile(tempDir)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestDataDecryptAES(t *testing.T) {
	key := make([]byte, 32) // AES-256 key
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("Failed to generate random key: %v", err)
	}

	plaintext := []byte("test plaintext")
	ciphertext, err := DataEncryptAES(plaintext, key)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	t.Run("ValidCiphertextAndKey", func(t *testing.T) {
		decrypted, err := DataDecryptAES(ciphertext, key)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if string(decrypted) != string(plaintext) {
			t.Errorf("Expected decrypted text to be %q, got %q", plaintext, decrypted)
		}
	})

	t.Run("ValidCiphertextAndIncorrectKey", func(t *testing.T) {
		incorrectKey := make([]byte, 32)
		if _, err := rand.Read(incorrectKey); err != nil {
			t.Fatalf("Failed to generate random incorrect key: %v", err)
		}
		_, err = DataDecryptAES(ciphertext, incorrectKey)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("ShortCiphertext", func(t *testing.T) {
		nonceSize := aes.BlockSize
		_, err := DataDecryptAES(ciphertext[:nonceSize-1], key)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("EmptyKey", func(t *testing.T) {
		_, err := DataDecryptAES(ciphertext, []byte{})
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("IncorrectKeyLength", func(t *testing.T) {
		incorrectKey := make([]byte, 16) // AES-128 key
		if _, err := rand.Read(incorrectKey); err != nil {
			t.Fatalf("Failed to generate random incorrect key: %v", err)
		}
		_, err = DataDecryptAES(ciphertext, incorrectKey)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestDataEncryptAES(t *testing.T) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		t.Fatalf("Failed to generate random key: %v", err)
	}
	plaintext := []byte("Hello, World!")
	ciphertext, err := DataEncryptAES(plaintext, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}
	if len(ciphertext) == 0 {
		t.Errorf("Ciphertext is empty")
	}

	emptyPlaintext := []byte{}
	_, err = DataEncryptAES(emptyPlaintext, key)
	if err != nil {
		t.Errorf("Encryption of empty plaintext failed: %v", err)
	}

	singleBytePlaintext := []byte{0x01}
	_, err = DataEncryptAES(singleBytePlaintext, key)
	if err != nil {
		t.Errorf("Encryption of single byte plaintext failed: %v", err)
	}

	_, err = DataEncryptAES(nil, key)
	if err != nil {
		t.Errorf("Encryption of nil plaintext failed: %v", err)
	}
}

func TestDataDecryptRSA(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key pair: %v", err)
	}

	plaintext := []byte("Hello, World!")
	ciphertext, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		&key.PublicKey,
		plaintext,
		nil)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decryptedText, err := DataDecryptRSA(ciphertext, key)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}
	if string(decryptedText) != string(plaintext) {
		t.Errorf("Decrypted text does not match original plaintext: %s != %s", decryptedText, plaintext)
	}

	emptyCiphertext := []byte{}
	_, err = DataDecryptRSA(emptyCiphertext, key)
	if err == nil {
		t.Errorf("Expected error for empty ciphertext, got nil")
	}

	singleByteCiphertext := []byte{0x01}
	_, err = DataDecryptRSA(singleByteCiphertext, key)
	if err == nil {
		t.Errorf("Expected error for single byte ciphertext, got nil")
	}

	_, err = DataDecryptRSA(nil, key)
	if err == nil {
		t.Errorf("Expected error for nil ciphertext, got nil")
	}

	invalidCiphertext := make([]byte, 256)
	if _, err := io.ReadFull(rand.Reader, invalidCiphertext); err != nil {
		t.Fatalf("Failed to generate random invalid ciphertext: %v", err)
	}
	_, err = DataDecryptRSA(invalidCiphertext, key)
	if err == nil {
		t.Errorf("Expected error for invalid ciphertext, got nil")
	}
}

func TestDataEncryptRSA(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key pair: %v", err)
	}

	plaintext := []byte("Hello, World!")
	ciphertext, err := DataEncryptRSA(plaintext, &key.PublicKey)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}
	if len(ciphertext) == 0 {
		t.Errorf("Ciphertext is empty")
	}

	decryptedText, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		key,
		ciphertext,
		nil)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}
	if string(decryptedText) != string(plaintext) {
		t.Errorf("Decrypted text does not match original plaintext: %s != %s", decryptedText, plaintext)
	}

	emptyPlaintext := []byte{}
	ciphertext, err = DataEncryptRSA(emptyPlaintext, &key.PublicKey)
	if err != nil {
		t.Errorf("Encryption of empty plaintext failed: %v", err)
	}
	if len(ciphertext) == 0 {
		t.Errorf("Ciphertext for empty plaintext is empty")
	}

	singleBytePlaintext := []byte{0x01}
	ciphertext, err = DataEncryptRSA(singleBytePlaintext, &key.PublicKey)
	if err != nil {
		t.Errorf("Encryption of single byte plaintext failed: %v", err)
	}
	if len(ciphertext) == 0 {
		t.Errorf("Ciphertext for single byte plaintext is empty")
	}

}
