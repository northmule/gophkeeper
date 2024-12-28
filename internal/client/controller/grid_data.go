package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/client/service"
	"github.com/northmule/gophkeeper/internal/common/model_data"
	"golang.org/x/net/context"
)

// GridData контроллер
type GridData struct {
	logger *logger.Logger
	cfg    *config.Config
	crypt  service.Cryptographer
}

// NewGridData конструктор
func NewGridData(cfg *config.Config, crypt service.Cryptographer, logger *logger.Logger) *GridData {
	return &GridData{
		logger: logger,
		cfg:    cfg,
		crypt:  crypt,
	}
}

// GridDataResponse ответ
type GridDataResponse struct {
	model_data.ListDataItemsResponse
}

// Send отправка запроса к серверу
func (c *GridData) Send(token string) (*GridDataResponse, error) {
	requestURL := fmt.Sprintf("%s/api/v1/items_list?offset=0&limit=200", c.cfg.Value().ServerAddress) //todo
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

	if response.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("вы не авторизованы")
	}

	if response.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("ошибка в запросе")
	}

	if response.StatusCode != http.StatusOK {
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

	responseData := new(GridDataResponse)
	err = json.Unmarshal(bodyRaw, responseData)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	return responseData, nil
}
