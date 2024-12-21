package keys

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"os"
	"path"
	"time"
)

const (
	// PrivateKeyFileName навзание файла ключа
	PrivateKeyFileName = "private_key.pem"
	// PublicKeyFileName навзание файла ключа
	PublicKeyFileName = "public_key.pem"
	// CertificateFileName навзание файла ключа
	CertificateFileName = "cert.pem"
	//PrivateKeyFileNameForEncryption Ключ для шифрования данных (есть на клиенте и на сервере)
	PrivateKeyFileNameForEncryption = "private_key_for_encryption.key"
)

// Keys сервис сертификата
type Keys struct {
	generator      KeyGenerator
	privateKeyPath string
	publicKeyPath  string
	certPath       string

	organization string
	country      string
	serialNumber *big.Int

	pubKey *rsa.PublicKey

	overwriting bool
}

// KeyGenerator интерфейс получения получение crypto.Signer
type KeyGenerator interface {
	GenerateKey() (crypto.Signer, error)
}

type Options struct {
	Generator    KeyGenerator
	SavePath     string
	Organization string
	Country      string
	// serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	SerialNumber *big.Int
}

// NewKeys конструктор
func NewKeys(options Options) *Keys {

	return &Keys{
		generator:      options.Generator,
		privateKeyPath: path.Join(options.SavePath, PrivateKeyFileName),
		publicKeyPath:  path.Join(options.SavePath, PublicKeyFileName),
		certPath:       path.Join(options.SavePath, CertificateFileName),
		organization:   options.Organization,
		country:        options.Country,
		serialNumber:   options.SerialNumber,
	}
}

// InitSelfSigned создаёт ключи и сертификат
func (c *Keys) InitSelfSigned() error {
	var err error
	privateKey, err := c.generator.GenerateKey()
	if err != nil {
		return err
	}
	if privateKey == nil {
		return errors.New("GenerateKey empty")
	}

	// шаблон сертификата
	certificateTemplate := x509.Certificate{
		SerialNumber: c.serialNumber,
		Subject: pkix.Name{
			Organization: []string{c.organization},
			Country:      []string{c.country},
		},
		// разрешаем использование сертификата для 127.0.0.1 и ::1
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:    []string{"localhost"},
		NotBefore:   time.Now(),
		// время жизни сертификата
		NotAfter: time.Now().Add(8760 * time.Hour),

		KeyUsage: x509.KeyUsageDigitalSignature,
		// Набобр для вариантов использования, example: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageOCSPSigning},
		BasicConstraintsValid: true,
	}

	// сертификат
	certBytes, err := x509.CreateCertificate(rand.Reader, &certificateTemplate, &certificateTemplate, privateKey.Public(), privateKey)
	if err != nil {
		return err
	}
	// кодирование сертификата в pem для хранения
	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	if pemCert == nil {
		return err
	}

	// Запись сертификата в файл
	if err = os.WriteFile(c.certPath, pemCert, 0644); err != nil {
		return err
	}

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return err
	}

	// кодирование закрытого ключа в pem для хранения
	pemPrivateKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyBytes})
	if pemPrivateKey == nil {
		return errors.New("failed to encode key to PEM")
	}
	// Запись ключа в файл
	if err = os.WriteFile(c.privateKeyPath, pemPrivateKey, 0644); err != nil {
		return err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return err
	}
	// кодирование открытого ключа в pem для хранения
	pemPubKey := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyBytes})
	if pemPubKey == nil {
		return errors.New("failed to encode key to PEM")
	}
	// Запись ключа в файл
	if err = os.WriteFile(c.publicKeyPath, pemPubKey, 0644); err != nil {
		return err
	}

	return nil
}

// CertPath путь к сертфификату
func (c *Keys) CertPath() string {
	return c.certPath
}

// PrivateKeyPath путь к ключу
func (c *Keys) PrivateKeyPath() string {
	return c.privateKeyPath
}

// PublicKeyPath путь к ключу
func (c *Keys) PublicKeyPath() string {
	return c.publicKeyPath
}
