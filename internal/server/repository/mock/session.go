package mock

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// MockSessionManager мок
type MockSessionManager struct {
	mock.Mock
}

// Add мок
func (m *MockSessionManager) Add(token string, expire time.Time) {
	m.Called(token, expire)
}

// IsValid мок
func (m *MockSessionManager) IsValid(token string) bool {
	args := m.Called(token)
	return args.Bool(0)
}
