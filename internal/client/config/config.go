package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	v     *viper.Viper
	value *ServerConfig
}

type ServerConfig struct {
	ServerAddress string `mapstructure:"ServerAddress"`
	LogLevel      string `mapstructure:"LogLever"`
	FilePath      string `mapstructure:"FilePath"`
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
