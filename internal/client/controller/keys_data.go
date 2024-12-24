package controller

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/client/service"
	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/keys"
	"github.com/northmule/gophkeeper/internal/common/util"
	"golang.org/x/net/context"
)

// KeysData контроллер обмена ключами с сервером
type KeysData struct {
	logger *logger.Logger
	cfg    *config.Config
	crypt  service.Cryptographer
}

// NewKeysData конструктор
func NewKeysData(cfg *config.Config, crypt service.Cryptographer, logger *logger.Logger) *KeysData {
	return &KeysData{
		logger: logger,
		cfg:    cfg,
		crypt:  crypt,
	}
}

// UploadClientPublicKey отправить публичный ключ на сервер
func (c *KeysData) UploadClientPublicKey(token string) error {
	requestURL := fmt.Sprintf("%s/api/v1/save_public_key", c.cfg.Value().ServerAddress)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(data_type.FileField, "public_key")
	if err != nil {
		return err
	}
	key, err := os.Open(path.Join(c.cfg.Value().PathKeys, keys.PublicKeyFileName))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, key)

	err = writer.Close()
	if err != nil {
		return err
	}
	ctx := context.Background()
	requestPrepare, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, body)
	if err != nil {
		return err
	}
	requestPrepare.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	requestPrepare.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	response, err := client.Do(requestPrepare)
	if err != nil {
		return err
	}
	if response.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("вы не авторизованы")
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("не известная ошибка")
	}
	defer response.Body.Close()

	return nil
}

// DownloadPublicServerKey забрать ключ с сервера
func (c *KeysData) DownloadPublicServerKey(token string) error {
	requestURL := fmt.Sprintf("%s/api/v1/download_server_public_key", c.cfg.Value().ServerAddress)
	ctx := context.Background()

	requestPrepare, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, nil)
	if err != nil {
		return err
	}
	requestPrepare.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	response, err := client.Do(requestPrepare)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		bodyRaw, _ := io.ReadAll(response.Body)
		return fmt.Errorf(string(bodyRaw))
	}

	if response.StatusCode == http.StatusBadRequest {
		bodyRaw, _ := io.ReadAll(response.Body)
		return fmt.Errorf(string(bodyRaw))
	}

	if response.StatusCode != http.StatusOK {
		bodyRaw, _ := io.ReadAll(response.Body)
		return fmt.Errorf(string(bodyRaw), response.StatusCode)
	}

	f, err := os.OpenFile(path.Join(c.cfg.Value().PathPublicKeyServer, keys.PublicKeyFileName), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, response.Body)
	if err != nil {
		return err
	}

	return nil
}

// UploadClientPrivateKey отправить приватный ключ ключ на сервер
func (c *KeysData) UploadClientPrivateKey(token string) error {
	requestURL := fmt.Sprintf("%s/api/v1/save_client_private_key", c.cfg.Value().ServerAddress)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(data_type.FileField, "private_client_key")
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path.Join(c.cfg.Value().PathKeys, keys.PrivateKeyFileNameForEncryption), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	privateKey := []byte(util.CreateHashForKey("super_secret_key")) //todo
	_, err = f.Write(privateKey)
	if err != nil {
		return err
	}

	// шифруем публичным серверным ключом
	privateKeyCrypt, err := c.crypt.EncryptRSA(privateKey)
	if err != nil {
		return err
	}
	_, err = part.Write(privateKeyCrypt)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}
	ctx := context.Background()
	requestPrepare, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, body)
	if err != nil {
		return err
	}
	requestPrepare.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	requestPrepare.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	response, err := client.Do(requestPrepare)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}
