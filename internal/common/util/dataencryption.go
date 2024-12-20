package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"os"
)

// DataEncryptRSA зашифровать данные публичным ключом
func DataEncryptRSA(rawDate []byte, key *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		key,
		rawDate,
		nil)
}

// DataDecryptRSA расшифровать данные
func DataDecryptRSA(encryptDate []byte, key *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		key,
		encryptDate,
		nil)
}

// DataEncryptAES зашифровать данные AES
func DataEncryptAES(rawDate []byte, key []byte) ([]byte, error) {
	block, _ := aes.NewCipher(key)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, rawDate, nil)
	return ciphertext, nil
}

// DataDecryptAES расшифровать данные
func DataDecryptAES(encryptDate []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := encryptDate[:nonceSize], encryptDate[nonceSize:]

	return gcm.Open(nil, nonce, ciphertext, nil)
}

// FillPrivateRsaKeyFromFile Вернёт приватный ключ RSA (для дешиврования входящих соообщениий)
func FillPrivateRsaKeyFromFile(path string) (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyBytes)

	if block == nil {
		return nil, fmt.Errorf("no PEM data found in file")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch key.(type) {
	case *rsa.PrivateKey:
		return key.(*rsa.PrivateKey), nil
	}

	return nil, fmt.Errorf("unsupported key type")
}

// FillPublicRsaKeyFromFile Вернёт публичный ключ RSA (для шифрования исходящих соообщениий)
func FillPublicRsaKeyFromFile(path string) (*rsa.PublicKey, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyBytes)

	if block == nil {
		return nil, fmt.Errorf("no PEM data found in file")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch key.(type) {
	case *rsa.PublicKey:
		return key.(*rsa.PublicKey), nil
	}
	return nil, fmt.Errorf("unsupported key type")
}

// FillPublicRsaKeyFromString Вернёт публичный ключ RSA (для шифрования исходящих соообщениий)
func FillPublicRsaKeyFromString(keyData string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(keyData))

	if block == nil {
		return nil, fmt.Errorf("no PEM data found in file")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch key.(type) {
	case *rsa.PublicKey:
		return key.(*rsa.PublicKey), nil
	}
	return nil, fmt.Errorf("unsupported key type")
}

// CreateHashForKey создаёт хэш для последующего использования в шифровании сообщений
func CreateHashForKey(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
