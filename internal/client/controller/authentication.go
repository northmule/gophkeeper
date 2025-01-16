package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"golang.org/x/net/context"
)

// Authentication контроллер
type Authentication struct {
	logger *logger.Logger
	cfg    *config.Config
}

// NewAuthentication конструктор
func NewAuthentication(cfg *config.Config, logger *logger.Logger) *Authentication {
	return &Authentication{
		logger: logger,
		cfg:    cfg,
	}
}

type authenticationRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// AuthenticationResponse ответ
type AuthenticationResponse struct {
	Value string
}

// Send отправка запроса к серверу
func (c *Authentication) Send(login string, password string) (*AuthenticationResponse, error) {
	requestURL := fmt.Sprintf("%s/api/v1/login", c.cfg.Value().ServerAddress)
	ctx := context.Background()

	requestData := authenticationRequest{
		Login:    login,
		Password: password,
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
		if response.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf("не верная пара логин/пароль")
		}

		return nil, fmt.Errorf("не известная ошибка")
	}

	var token string
	token = response.Header.Get("Authorization")
	token = strings.Replace(token, "Bearer ", "", 1)
	token = strings.Trim(token, " ")

	responseData := new(AuthenticationResponse)
	responseData.Value = token

	return responseData, nil
}
