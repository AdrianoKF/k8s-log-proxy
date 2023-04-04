package config

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	AllowedNamespaces []string
}

const CONFIG_NAME = "config"

// ReadConfig reads and parses the application config from a `config.*` file found in either the current working directory or the root directory
func ReadConfig() (appConfig AppConfig, err error) {
	viper.SetConfigName(CONFIG_NAME)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/")
	if err = viper.ReadInConfig(); err != nil {
		return
	}

	if err = viper.Unmarshal(&appConfig); err != nil {
		return
	}
	return
}
