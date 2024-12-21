package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"golang.org/x/net/context"
)

// Registration контроллер
type Registration struct {
	logger *logger.Logger
	cfg    *config.Config
}

// NewRegistration конструктор
func NewRegistration(cfg *config.Config, logger *logger.Logger) *Registration {
	return &Registration{
		logger: logger,
		cfg:    cfg,
	}
}

type registrationRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RegistrationResponse struct {
	Value string
}

// Send отправка запроса к серверу
func (c *Registration) Send(login string, password string, email string) (*RegistrationResponse, error) {
	requestURL := fmt.Sprintf("%s/api/v1/register", c.cfg.Value().ServerAddress)
	ctx := context.Background()

	requestData := registrationRequest{
		Login:    login,
		Password: password,
		Email:    email,
	}

	requestBody, err := json.Marshal(requestData)
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
	client := &http.Client{}
	response, err := client.Do(requestPrepare)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("не известная ошибка")
	}

	responseData := new(RegistrationResponse)
	responseData.Value = "ok"

	return responseData, nil
}
