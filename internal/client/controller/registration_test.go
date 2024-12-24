package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/stretchr/testify/assert"
)

func TestNewRegistration_Success(t *testing.T) {
	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}

	mockConfig := makeMockConfig("")

	controller := NewRegistration(mockConfig, log)
	assert.NotNil(t, controller)
	assert.Equal(t, mockConfig, controller.cfg)
	assert.Equal(t, log, controller.logger)
}

func TestRegistration_Send_Success(t *testing.T) {
	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/register", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var requestData registrationRequest
		err := json.NewDecoder(r.Body).Decode(&requestData)
		assert.NoError(t, err)
		assert.Equal(t, "testlogin", requestData.Login)
		assert.Equal(t, "testpassword", requestData.Password)
		assert.Equal(t, "testemail", requestData.Email)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Value": "ok"}`))
	}))
	defer testServer.Close()
	mockConfig := makeMockConfig(testServer.URL)
	controller := NewRegistration(mockConfig, log)
	response, err := controller.Send("testlogin", "testpassword", "testemail")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "ok", response.Value)
}

func TestRegistration_Send_HTTPRequestError(t *testing.T) {
	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}

	mockConfig := makeMockConfig("")
	controller := NewRegistration(mockConfig, log)

	response, err := controller.Send("testlogin", "testpassword", "testemail")
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestRegistration_Send_NonOKResponse(t *testing.T) {
	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Value": "error"}`))
	}))
	defer testServer.Close()
	mockConfig := makeMockConfig(testServer.URL)
	controller := NewRegistration(mockConfig, log)

	response, err := controller.Send("testlogin", "testpassword", "testemail")
	assert.Error(t, err)
	assert.Nil(t, response)
}
