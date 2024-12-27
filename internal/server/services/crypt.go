package service

import (
	"crypto/rsa"
	"path"

	"github.com/northmule/gophkeeper/internal/common/keys"
	"github.com/northmule/gophkeeper/internal/common/util"
	"github.com/northmule/gophkeeper/internal/server/config"
)

// Crypt сервис шифрования дешифрования данных
type Crypt struct {
	// Клиентский публичный ключ, для шифрования исходящиз сообщений
	clientPublicKey *rsa.PublicKey
	// Приватный ключ сервера, для расшифровки сообщений от клиента
	serverPrivateKey *rsa.PrivateKey

	cfg *config.Config
}

type CryptService interface {
	EncryptRSA(data []byte) ([]byte, error)
	DecryptRSA(data []byte) ([]byte, error)
}

// NewCrypt конструктор
func NewCrypt(cfg *config.Config) (*Crypt, error) {
	instance := new(Crypt)
	var err error
	instance.clientPublicKey, err = util.FillPublicRsaKeyFromFile(path.Join(cfg.Value().PathKeys, keys.PublicKeyFileName))
	if err != nil {
		return nil, err
	}
	instance.serverPrivateKey, err = util.FillPrivateRsaKeyFromFile(path.Join(cfg.Value().PathKeys, keys.PrivateKeyFileName))
	if err != nil {
		return nil, err
	}
	instance.cfg = cfg
	return instance, nil
}

// EncryptRSA Шифрование исходящих данных клиентским публичным ключом
func (crypt *Crypt) EncryptRSA(data []byte) ([]byte, error) {
	return util.DataEncryptRSA(data, crypt.clientPublicKey)
}

// DecryptRSA Расшифровка входящих сообещний приватным ключом сервера
func (crypt *Crypt) DecryptRSA(data []byte) ([]byte, error) {
	return util.DataDecryptRSA(data, crypt.serverPrivateKey)
}
