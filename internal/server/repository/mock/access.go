package mock

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

type MockAccessService struct {
	mock.Mock
}

func (m *MockAccessService) PasswordHash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}
func (m *MockAccessService) FillJWTToken() *jwtauth.JWTAuth {
	args := m.Called()
	return args.Get(0).(*jwtauth.JWTAuth)
}

func (m *MockAccessService) GetUserUUIDByJWTToken(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockAccessService) FindTokenByRequest(r *http.Request) string {
	args := m.Called(r)
	return args.String(0)
}
