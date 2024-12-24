package controller

import (
	"bytes"
	"encoding/json"
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
	"github.com/northmule/gophkeeper/internal/common/model_data"
	"golang.org/x/net/context"
)

// FileData контроллер
type FileData struct {
	logger *logger.Logger
	cfg    *config.Config
	crypt  service.Cryptographer
}

// NewFileData конструктор
func NewFileData(cfg *config.Config, crypt service.Cryptographer, logger *logger.Logger) *FileData {
	return &FileData{
		logger: logger,
		cfg:    cfg,
		crypt:  crypt,
	}
}

type FileDataResponse struct {
	// Адрес без хоста для загрузки данных файла
	UploadPath string `json:"upload_path"`
}

// Send отправка запроса к серверу. Предзагрузка основной информации о файле. В ответе будет адрес куда отправлять сам файл
func (c *FileData) Send(token string, requestData *model_data.FileDataInitRequest) (*FileDataResponse, error) {
	requestURL := fmt.Sprintf("%s/api/v1/file_data/init", c.cfg.Value().ServerAddress)
	ctx := context.Background()

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}
	// Шифруем
	requestBody, err = c.crypt.EncryptAES(requestBody)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	buf := bytes.NewBuffer(requestBody)
	requestPrepare, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, buf)
	if err != nil {
		return nil, err
	}
	requestPrepare.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	response, err := client.Do(requestPrepare)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bodyRaw, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("вы не авторизованы")
	}

	if response.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("ошибка в запросе")
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("не известная ошибка")
	}

	responseData := new(FileDataResponse)

	err = json.Unmarshal(bodyRaw, responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

// UploadFile отправка файла на сервер
func (c *FileData) UploadFile(token string, url string, file *os.File) error {
	requestURL := fmt.Sprintf("%s/api/v1%s", c.cfg.Value().ServerAddress, url)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(data_type.FileField, file.Name())
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return err
	}

	var buf []byte
	// Шифруем файл
	buf, err = c.crypt.EncryptAES(body.Bytes())
	if err != nil {
		c.logger.Error(err)
		return err
	}
	bodyCrypt := &bytes.Buffer{}
	_, err = bodyCrypt.Write(buf)
	if err != nil {
		c.logger.Error(err)
		return err
	}
	ctx := context.Background()
	requestPrepare, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bodyCrypt)
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
	defer response.Body.Close()

	return nil
}

func (c *FileData) DownLoadFile(token string, fileName string, dataUUID string) error {
	requestURL := fmt.Sprintf("%s/api/v1/file_data/get/%s/%s", c.cfg.Value().ServerAddress, dataUUID, "0")
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

	f, err := os.OpenFile(path.Join(c.cfg.Value().FilePath, fileName), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	bodyRaw, err := io.ReadAll(response.Body)
	if err != nil {
		c.logger.Error(err)
		return err
	}
	// Расшифровка тела
	bodyRaw, err = c.crypt.DecryptAES(bodyRaw)
	if err != nil {
		c.logger.Error(err)
		return err
	}
	_, err = f.Write(bodyRaw)
	if err != nil {
		return err
	}

	return nil
}
