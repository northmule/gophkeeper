package access

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/northmule/gophkeeper/internal/common/util"
	"github.com/northmule/gophkeeper/internal/server/api/rctx"
	"github.com/northmule/gophkeeper/internal/server/config"
)

const CookieAuthName = "gophermart_session"
const HMACTokenExp = time.Hour * 600
const HMACSecretKey = "super_secret_key_gophermart"

type Access struct {
	cfg *config.Config
}

func NewAccess(cfg *config.Config) *Access {
	return &Access{cfg: cfg}
}

func (a *Access) PasswordHash(password string) (string, error) {
	switch a.cfg.Value().PasswordAlgoHashing {
	case "sha256":
		return util.PasswordHashSha256(password), nil
	case "sha512":
		return util.PasswordHashSha512(password), nil
	}
	return "", fmt.Errorf("unknown hashing algorithm")
}

func (a *Access) FillJWTToken() *jwtauth.JWTAuth {
	return jwtauth.New("HS512", []byte("_secret_"), nil) // todo в конфиг
}

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

func (a *Access) FindTokenByRequest(r *http.Request) string {
	var token string
	token = r.Header.Get("Authorization")
	token = strings.Replace(token, "Bearer ", "", 1)
	token = strings.Trim(token, " ")
	return token
}

func (a *Access) GetUserToken(req *http.Request) (string, error) {
	token := req.Header.Get("Authorization")
	if token == "" {
		cookieAuth, err := req.Cookie(CookieAuthName)
		if err != nil {
			return "", err
		}
		token = cookieAuth.Value
	}

	return token, nil
}

func (a *Access) GenerateToken(userUUID string, exp time.Duration, secretKey string) (string, time.Time) {
	hashed := hmac.New(sha256.New, []byte(secretKey))
	hashed.Write([]byte(userUUID))
	token := hex.EncodeToString(hashed.Sum(nil))
	tokenExp := time.Now().Add(exp)
	return token, tokenExp
}

func (a *Access) ValidateToken(userUUID string, token string, secretKey string) (bool, error) {
	tokenSign, err := hex.DecodeString(token)
	if err != nil {
		return false, err
	}
	hashed := hmac.New(sha256.New, []byte(secretKey))
	hashed.Write([]byte(userUUID))
	expectedSign := hashed.Sum(nil)

	if !hmac.Equal(tokenSign, expectedSign) {
		return false, err
	}
	return true, nil
}

func (a *Access) Authentication(userUUID string) (string, string, *http.Cookie, *time.Time) {

	token, tokenExp := a.GenerateToken(userUUID, HMACTokenExp, HMACSecretKey)
	tokenValue := fmt.Sprintf("%s:%s", token, userUUID)

	cookie := &http.Cookie{
		Name:    CookieAuthName,
		Value:   tokenValue,
		Expires: tokenExp,
		Secure:  false,
		Path:    "/",
	}

	return token, tokenValue, cookie, &tokenExp
}
