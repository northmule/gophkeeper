package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/client/service"
	"github.com/northmule/gophkeeper/internal/common/model_data"
	"golang.org/x/net/context"
)

// CardData контроллер
type CardData struct {
	logger *logger.Logger
	cfg    *config.Config
	crypt  service.Cryptographer
}

// NewCardData конструктор
func NewCardData(cfg *config.Config, crypt service.Cryptographer, logger *logger.Logger) *CardData {
	return &CardData{
		logger: logger,
		cfg:    cfg,
		crypt:  crypt,
	}
}

// CardDataResponse ответ
type CardDataResponse struct {
	Value string
}

// Send отправка запроса к серверу
func (c *CardData) Send(token string, requestData *model_data.CardDataRequest) (*CardDataResponse, error) {
	requestURL := fmt.Sprintf("%s/api/v1/save_card_data", c.cfg.Value().ServerAddress)
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
		c.logger.Error(err)
		return nil, err
	}
	requestPrepare.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	response, err := client.Do(requestPrepare)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("вы не авторизованы")
	}

	if response.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("ошибка в запросе")
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("не известная ошибка")
	}

	responseData := new(CardDataResponse)
	responseData.Value = "ok"

	return responseData, nil
}
