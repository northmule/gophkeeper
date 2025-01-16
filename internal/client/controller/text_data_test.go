package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/common/model_data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCryptographer мок
type MockCryptographer struct {
	mock.Mock
}

func (m *MockCryptographer) EncryptAES(data []byte) ([]byte, error) {
	args := m.Called(data)

	if _, ok := (args.Get(0)).([]byte); ok {
		return (args.Get(0)).([]byte), args.Error(1)
	} else {
		return []byte(""), args.Error(1)
	}

}

func (m *MockCryptographer) DecryptAES(data []byte) ([]byte, error) {
	args := m.Called(data)
	return args.Get(0).([]byte), args.Error(1)
}

// EncryptRSA Шифрование исходящих данных серверных публичным ключом
func (m *MockCryptographer) EncryptRSA(data []byte) ([]byte, error) {
	return nil, nil
}

// DecryptRSA Расшифровка входящих сообещний приватным ключом клиента
func (m *MockCryptographer) DecryptRSA(data []byte) ([]byte, error) {
	return nil, nil
}

func TestNewTextData_Success(t *testing.T) {
	mockConfig := makeMockConfig("")
	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	mockCrypt := new(MockCryptographer)

	textData := NewTextData(mockConfig, mockCrypt, log)
	assert.NotNil(t, textData)
	assert.Equal(t, mockConfig, textData.cfg)
	assert.Equal(t, mockCrypt, textData.crypt)
	assert.Equal(t, log, textData.logger)
}

func TestTextData_Send_Success(t *testing.T) {

	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	mockCrypt := new(MockCryptographer)

	mockCrypt.On("EncryptAES", mock.Anything).Return([]byte("encrypted_data"), nil)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/save_text_data", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "Bearer test_token", r.Header.Get("Authorization"))

		var requestData model_data.TextDataRequest
		err = json.Unmarshal([]byte(`{"name":"name"}`), &requestData)
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Value": "ok"}`))
	}))
	defer testServer.Close()
	mockConfig := makeMockConfig(testServer.URL)
	textData := NewTextData(mockConfig, mockCrypt, log)
	requestData := &model_data.TextDataRequest{
		Value: "test_text",
	}
	response, err := textData.Send("test_token", requestData)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "ok", response.Value)
}

func TestTextData_Send_HTTPRequestError(t *testing.T) {

	mockConfig := makeMockConfig("")
	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	mockCrypt := new(MockCryptographer)
	mockCrypt.On("EncryptAES", mock.Anything).Return([]byte("encrypted_data"), nil)

	textData := NewTextData(mockConfig, mockCrypt, log)
	requestData := &model_data.TextDataRequest{
		Value: "test_text",
	}
	response, err := textData.Send("test_token", requestData)
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestTextData_Send_EncryptionError(t *testing.T) {
	mockConfig := makeMockConfig("")
	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	mockCrypt := new(MockCryptographer)

	mockCrypt.On("EncryptAES", mock.Anything).Return(nil, errors.New("encryption error"))

	textData := NewTextData(mockConfig, mockCrypt, log)
	requestData := &model_data.TextDataRequest{
		Value: "test_text",
	}
	response, err := textData.Send("test_token", requestData)
	assert.Error(t, err)
	assert.Nil(t, response)
}
