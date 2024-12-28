package access

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/jwtauth/v5"
	"github.com/northmule/gophkeeper/internal/common/util"
	"github.com/northmule/gophkeeper/internal/server/api/rctx"
	"github.com/northmule/gophkeeper/internal/server/config"
)

// Access сервис проверки доступа
type Access struct {
	cfg *config.Config
}

// AccessService серис доступа
type AccessService interface {
	PasswordHash(password string) (string, error)
	FillJWTToken() *jwtauth.JWTAuth
	GetUserUUIDByJWTToken(ctx context.Context) (string, error)
	FindTokenByRequest(r *http.Request) string
}

// NewAccess конструктор
func NewAccess(cfg *config.Config) *Access {
	return &Access{cfg: cfg}
}

// PasswordHash хэшер пароля
func (a *Access) PasswordHash(password string) (string, error) {
	switch a.cfg.Value().PasswordAlgoHashing {
	case "sha256":
		return util.PasswordHashSha256(password), nil
	case "sha512":
		return util.PasswordHashSha512(password), nil
	}
	return "", fmt.Errorf("unknown hashing algorithm")
}

// FillJWTToken начальное значение токена
func (a *Access) FillJWTToken() *jwtauth.JWTAuth {
	return jwtauth.New("HS512", []byte("_secret_"), nil) // todo в конфиг
}

// GetUserUUIDByJWTToken UUID пользвоателя из токена
func (a *Access) GetUserUUIDByJWTToken(ctx context.Context) (string, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return "", err
	}
	value, ok := claims[rctx.MapKeyUserUUID]
	if !ok {
		return "", fmt.Errorf("no uuid found in claims")
	}
	return value.(string), nil
}

// FindTokenByRequest поиск токена в запросе
func (a *Access) FindTokenByRequest(r *http.Request) string {
	var token string
	token = r.Header.Get("Authorization")
	token = strings.Replace(token, "Bearer ", "", 1)
	token = strings.Trim(token, " ")
	return token
}
