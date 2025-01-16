package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigInit(t *testing.T) {
	t.Run("Valid .env file", func(t *testing.T) {
		validEnvContent := `ADDRESS=localhost:8080
DSN=postgres://user:pass@localhost/db
LOG_LEVEL=info
HTTP_COMPRESS_LEVEL=5
PASSWORD_ALGO_HASHING=bcrypt
PATH_FILE_STORAGE=/var/files
PATH_KEYS=/var/keys
OVERWRITE_KEYS=true`

		validConfigPath := filepath.Join(".server.env")
		if err := os.WriteFile(validConfigPath, []byte(validEnvContent), 0644); err != nil {
			t.Fatalf("Failed to write valid config file: %v", err)
		}
		defer os.Remove(validConfigPath)
		validEnvFile, err := os.CreateTemp(t.TempDir(), "valid.env")
		require.NoError(t, err)
		defer os.Remove(validEnvFile.Name())
		_, err = validEnvFile.WriteString(validEnvContent)
		require.NoError(t, err)
		validEnvFile.Close()
		cfg := NewConfig()
		err = cfg.Init()
		assert.NoError(t, err)
		serverConfig := cfg.Value()

		wantValidConfig := &ServerConfig{
			Address:             "localhost:8080",
			Dsn:                 "postgres://user:pass@localhost/db",
			LogLevel:            "info",
			HTTPCompressLevel:   5,
			PasswordAlgoHashing: "bcrypt",
			PathFileStorage:     "/var/files",
			PathKeys:            "/var/keys",
			OverwriteKeys:       true,
		}
		if diff := cmp.Diff(wantValidConfig, serverConfig); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}

	})

}
