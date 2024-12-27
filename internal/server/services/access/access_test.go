package access

import (
	"net/http"
	"testing"

	"github.com/northmule/gophkeeper/internal/common/util"
	"github.com/northmule/gophkeeper/internal/server/config"
	"github.com/stretchr/testify/assert"
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
