package mock

import "github.com/stretchr/testify/mock"

// MockCryptService мок
type MockCryptService struct {
	mock.Mock
}

// EncryptRSA мок
func (m *MockCryptService) EncryptRSA(data []byte) ([]byte, error) {
	args := m.Called(data)
	return args.Get(0).([]byte), args.Error(1)
}

// DecryptRSA мок
func (m *MockCryptService) DecryptRSA(data []byte) ([]byte, error) {
	args := m.Called(data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}
