package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/client/service"
	"github.com/northmule/gophkeeper/internal/common/model_data"
	"golang.org/x/net/context"
)

// ItemData контроллер запроса данных по uuid
type ItemData struct {
	logger *logger.Logger
	cfg    *config.Config
	crypt  service.Cryptographer
}

// NewItemData конструктор
func NewItemData(cfg *config.Config, crypt service.Cryptographer, logger *logger.Logger) *ItemData {
	return &ItemData{
		logger: logger,
		cfg:    cfg,
		crypt:  crypt,
	}
}

// Send отправка запроса к серверу
func (c *ItemData) Send(token string, dataUUID string) (*model_data.DataByUUIDResponse, error) {
	var requestURL string
	requestURL = fmt.Sprintf("%s/api/v1/item_get/{uuid}", c.cfg.Value().ServerAddress)
	requestURL = strings.Replace(requestURL, "{uuid}", dataUUID, 1)
	ctx := context.Background()

	requestPrepare, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
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

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf("вы не авторизованы")
		}

		if response.StatusCode == http.StatusBadRequest {
			return nil, fmt.Errorf("ошибка в запросе")
		}
		return nil, fmt.Errorf("не известная ошибка")
	}

	bodyRaw, err := io.ReadAll(response.Body)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	// Расшифровка тела
	bodyRaw, err = c.crypt.DecryptAES(bodyRaw)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	responseData := new(model_data.DataByUUIDResponse)
	err = json.Unmarshal(bodyRaw, responseData)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	return responseData, nil
}
