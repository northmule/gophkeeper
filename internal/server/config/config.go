package config

import (
	"github.com/spf13/viper"
)

// Config конфигуратор
type Config struct {
	v     *viper.Viper
	value *ServerConfig
}

// ServerConfig конфигурация сервера
type ServerConfig struct {
	Address             string `mapstructure:"ADDRESS"`
	Dsn                 string `mapstructure:"DSN"`
	LogLevel            string `mapstructure:"LOG_LEVEL"`
	HTTPCompressLevel   int    `mapstructure:"HTTP_COMPRESS_LEVEL"`
	PasswordAlgoHashing string `mapstructure:"PASSWORD_ALGO_HASHING"`
	// PathFileStorage файлы пользователей
	PathFileStorage string `mapstructure:"PATH_FILE_STORAGE"`
	// PathKeys место хранения ключей сервера
	PathKeys      string `mapstructure:"PATH_KEYS"`
	OverwriteKeys bool   `mapstructure:"OVERWRITE_KEYS"`
}

// ErrorCfg сообщение с ошибкой
type ErrorCfg error

// NewConfig конструктор
func NewConfig() *Config {
	instance := new(Config)
	instance.v = viper.New()
	instance.value = new(ServerConfig)
	return instance
}

// Init собирает конфиг
func (c *Config) Init() error {
	var err error
	c.v.AddConfigPath(".")
	c.v.SetConfigName(".server")
	c.v.SetConfigType("env")
	err = c.v.ReadInConfig()
	if err != nil {
		return ErrorCfg(err)
	}

	err = c.v.Unmarshal(c.value)
	if err != nil {
		return ErrorCfg(err)
	}

	return nil
}

// Value значения конфига
func (c *Config) Value() *ServerConfig {
	return c.value
}
