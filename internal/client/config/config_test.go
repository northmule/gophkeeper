package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInit(t *testing.T) {
	testDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	validConfig := `
ServerAddress: "localhost:8080"
LogLevel: "info"
FilePath: "/path/to/file"
PathKeys: "/path/to/keys"
PathPublicKeyServer: "/path/to/server/public/key"
OverwriteKeys: true
`

	invalidConfig := `
ServerAddress: "localhost:8080"
LogLevel: "info"
FilePath: "/path/to/file"
PathKeys: "/path/to/keys"
PathPublicKeyServer: "/path/to/server/public/key"
OverwriteKeys: "true"  # Invalid type
`

	missingFieldConfig := `
ServerAddress: "localhost:8080"
LogLevel: "info"
FilePath: "/path/to/file"
PathKeys: "/path/to/keys"
# Missing PathPublicKeyServer
OverwriteKeys: true
`

	validConfigPath := filepath.Join("client.yaml")
	invalidConfigPath := filepath.Join("client.yaml")
	missingFieldConfigPath := filepath.Join("client.yaml")

	if err := os.WriteFile(validConfigPath, []byte(validConfig), 0644); err != nil {
		t.Fatalf("Failed to write valid config file: %v", err)
	}

	defer os.Remove(validConfigPath)
	defer os.Remove(invalidConfigPath)
	defer os.Remove(missingFieldConfigPath)

	cfg := NewConfig()
	cfg.v.SetConfigFile(validConfigPath)
	err = cfg.Init()
	if err != nil {
		t.Errorf("Init failed with valid config: %v", err)
	}

	value := cfg.Value()

	wantValidConfig := &ServerConfig{
		ServerAddress:       "localhost:8080",
		LogLevel:            "info",
		FilePath:            "/path/to/file",
		PathKeys:            "/path/to/keys",
		PathPublicKeyServer: "/path/to/server/public/key",
		OverwriteKeys:       true,
	}
	if diff := cmp.Diff(wantValidConfig, value); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}

	if err := os.WriteFile(invalidConfigPath, []byte(invalidConfig), 0644); err != nil {
		t.Fatalf("Failed to write invalid config file: %v", err)
	}

	cfg2 := NewConfig()
	cfg2.v.SetConfigFile(invalidConfigPath)
	err = cfg2.Init()
	if err != nil {
		t.Errorf(err.Error())
	}

	if err := os.WriteFile(missingFieldConfigPath, []byte(missingFieldConfig), 0644); err != nil {
		t.Fatalf("Failed to write missing field config file: %v", err)
	}

	cfg3 := NewConfig()
	cfg3.v.SetConfigFile(missingFieldConfigPath)
	err = cfg3.Init()
	if err != nil {
		t.Errorf(err.Error())
	}
}
