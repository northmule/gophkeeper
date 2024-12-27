package access

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt/v4"
	"github.com/northmule/gophkeeper/internal/common/util"
	"github.com/northmule/gophkeeper/internal/server/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestAccess_PasswordHash(t *testing.T) {
	cfg := config.NewConfig()
	_ = cfg.Init()
	cfg.Value().PasswordAlgoHashing = "sha256"

	access := NewAccess(cfg)
	hash, err := access.PasswordHash("password")
	assert.NoError(t, err)
	assert.Equal(t, util.PasswordHashSha256("password"), hash)

	cfg.Value().PasswordAlgoHashing = "sha512"

	hash, err = access.PasswordHash("password")
	assert.NoError(t, err)
	assert.Equal(t, util.PasswordHashSha512("password"), hash)

	cfg.Value().PasswordAlgoHashing = "unknown"

	hash, err = access.PasswordHash("password")
	assert.Error(t, err)
	assert.Empty(t, hash)

}

func TestAccess_FillJWTToken(t *testing.T) {
	cfg := config.NewConfig()
	_ = cfg.Init()

	access := NewAccess(cfg)
	token := access.FillJWTToken()
	assert.NotNil(t, token)
}

func TestAccess_FindTokenByRequest(t *testing.T) {
	cfg := config.NewConfig()
	_ = cfg.Init()

	access := NewAccess(cfg)
	req := &http.Request{
		Header: http.Header{"Authorization": {"Bearer test-token"}},
	}
	token := access.FindTokenByRequest(req)
	assert.Equal(t, "test-token", token)

	req = &http.Request{
		Header: http.Header{"Authorization": {"test-token"}},
	}
	token = access.FindTokenByRequest(req)
	assert.Equal(t, "test-token", token)

	req = &http.Request{
		Header: http.Header{},
	}
	token = access.FindTokenByRequest(req)
	assert.Equal(t, "", token)
}

//

type MockJWTAuth struct {
	mock.Mock
}

func (m *MockJWTAuth) FromContext(ctx context.Context) (string, map[string]interface{}, error) {
	args := m.Called(ctx)
	return args.String(0), args.Get(1).(map[string]interface{}), args.Error(2)
}

type jWTClaims struct {
	jwt.RegisteredClaims
	UserUUID string `json:"user_uuid"`
}

type jWTClaimsBad struct {
	jwt.RegisteredClaims
	User string `json:"user"`
}

func TestGetUserUUIDByJWTToken(t *testing.T) {
	cfg := config.NewConfig()
	_ = cfg.Init()

	t.Run("Valid JWT Token with UUID", func(t *testing.T) {
		ctx := context.Background()
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, jWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 300)),
			},
			UserUUID: "123e4567-e89b-12d3-a456-426614174000",
		})
		tokenValue, err := token.SignedString([]byte("_secret_"))
		jj := jwtauth.New("HS512", []byte("_secret_"), nil)
		jwtV, _ := jj.Decode(tokenValue)
		ctx = jwtauth.NewContext(ctx, jwtV, nil)

		a := NewAccess(cfg)
		uuid, err := a.GetUserUUIDByJWTToken(ctx)
		require.NoError(t, err)
		assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", uuid)
	})
	t.Run("without UUID", func(t *testing.T) {
		ctx := context.Background()
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, jWTClaimsBad{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 300)),
			},
		})
		tokenValue, err := token.SignedString([]byte("_secret_"))
		jj := jwtauth.New("HS512", []byte("_secret_"), nil)
		jwtV, _ := jj.Decode(tokenValue)
		ctx = jwtauth.NewContext(ctx, jwtV, nil)

		a := NewAccess(cfg)
		uuid, err := a.GetUserUUIDByJWTToken(ctx)
		require.Error(t, err)
		assert.Empty(t, uuid)
	})

}
