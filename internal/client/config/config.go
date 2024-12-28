package config

import (
	"github.com/spf13/viper"
)

// Config конфигурация клиента
type Config struct {
	v     *viper.Viper
	value *ServerConfig
}

// ServerConfig структура конфигурации клиента
type ServerConfig struct {
	ServerAddress string `mapstructure:"ServerAddress"`
	LogLevel      string `mapstructure:"LogLevel"`
	FilePath      string `mapstructure:"FilePath"`
	// Путь для сохранения публичного и приватного ключей клиента
	PathKeys string `json:"PathKeys"`
	// Путь к папке с публичным ключем сервера
	PathPublicKeyServer string `json:"PathPublicKeyServer"`
	// Перезаписывать клиентские ключи при старте клиента
	OverwriteKeys bool `json:"OverwriteKeys"`
}

// ErrorCfg ошибка конфигурации
type ErrorCfg error

// NewConfig конструктор
func NewConfig() *Config {
	instance := new(Config)
	instance.v = viper.New()
	instance.value = new(ServerConfig)
	return instance
}

// Init инициализация настроек из файла
func (c *Config) Init() error {
	var err error
	c.v.AddConfigPath(".")
	c.v.SetConfigName("client")
	c.v.SetConfigType("yaml")
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

func (c *Config) Value() *ServerConfig {
	return c.value
}
