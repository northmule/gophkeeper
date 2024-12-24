package service

import (
	"crypto/rsa"
	"os"
	"path"

	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/common/keys"
	"github.com/northmule/gophkeeper/internal/common/util"
)

// Crypt сервис шифрования дешифрования данных
type Crypt struct {
	// Серверный публичный ключ, для шифрования исходящиз сообщений
	serverPublicKey *rsa.PublicKey
	// Приватный ключ этого клиента, для расшифровки сообщений от сервера
	clientPrivateKey *rsa.PrivateKey
	// Ключ для шифрования и дешифрования данных между клиентом и сервером (ключ хранится и на клиенте и на сервере)
	privateKeyForEncryption []byte

	cfg *config.Config
}

// Cryptographer общий интерфейс шифрования для клиента
type Cryptographer interface {
	// EncryptRSA Шифрование исходящих данных серверных публичным ключом
	EncryptRSA(data []byte) ([]byte, error)
	// DecryptRSA Расшифровка входящих сообещний приватным ключом клиента
	DecryptRSA(data []byte) ([]byte, error)
	// EncryptAES Шифрование исходящих данных
	EncryptAES(data []byte) ([]byte, error)
	// DecryptAES Расшифровка входящих сообещний
	DecryptAES(data []byte) ([]byte, error)
}

// NewCrypt конструктор
func NewCrypt(cfg *config.Config) (*Crypt, error) {
	instance := new(Crypt)
	var err error
	instance.serverPublicKey, err = util.FillPublicRsaKeyFromFile(path.Join(cfg.Value().PathPublicKeyServer, keys.PublicKeyFileName))
	if err != nil {
		return nil, err
	}
	instance.clientPrivateKey, err = util.FillPrivateRsaKeyFromFile(path.Join(cfg.Value().PathKeys, keys.PrivateKeyFileName))
	if err != nil {
		return nil, err
	}
	instance.privateKeyForEncryption, err = os.ReadFile(path.Join(cfg.Value().PathKeys, keys.PrivateKeyFileNameForEncryption))
	if err != nil {
		return nil, err
	}
	instance.cfg = cfg
	return instance, nil
}

// EncryptRSA Шифрование исходящих данных серверных публичным ключом
func (crypt *Crypt) EncryptRSA(data []byte) ([]byte, error) {
	return util.DataEncryptRSA(data, crypt.serverPublicKey)
}

// DecryptRSA Расшифровка входящих сообещний приватным ключом клиента
func (crypt *Crypt) DecryptRSA(data []byte) ([]byte, error) {
	return util.DataDecryptRSA(data, crypt.clientPrivateKey)
}

// EncryptAES Шифрование исходящих данных
func (crypt *Crypt) EncryptAES(data []byte) ([]byte, error) {
	return util.DataEncryptAES(data, crypt.privateKeyForEncryption)
}

// DecryptAES Расшифровка входящих сообещний
func (crypt *Crypt) DecryptAES(data []byte) ([]byte, error) {
	return util.DataDecryptAES(data, crypt.privateKeyForEncryption)
}
