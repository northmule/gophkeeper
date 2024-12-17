package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	v     *viper.Viper
	value *ServerConfig
}

type ServerConfig struct {
	Address             string `mapstructure:"ADDRESS"`
	Dsn                 string `mapstructure:"DSN"`
	LogLevel            string `mapstructure:"LOG_LEVEL"`
	HTTPCompressLevel   int    `mapstructure:"HTTP_COMPRESS_LEVEL"`
	PasswordAlgoHashing string `mapstructure:"PASSWORD_ALGO_HASHING"`
	PathFileStorage     string `mapstructure:"PATH_FILE_STORAGE"`
}

type ErrorCfg error

func NewConfig() *Config {
	instance := new(Config)
	instance.v = viper.New()
	instance.value = new(ServerConfig)
	return instance
}

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

func (c *Config) Value() *ServerConfig {
	return c.value
}
