package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type AppConfig struct {
	AllowedNamespaces []string
}

const CONFIG_NAME = "config"

var Config AppConfig

func Init() {
	viper.SetConfigName(CONFIG_NAME)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		Config, _ = ReadConfig()
	})
	var err error
	if Config, err = ReadConfig(); err != nil {
		panic(err.Error())
	}
}

// ReadConfig reads and parses the application config from a `config.*` file found in either the current working directory or the root directory
func ReadConfig() (appConfig AppConfig, err error) {
	if err = viper.ReadInConfig(); err != nil {
		return
	}

	if err = viper.Unmarshal(&appConfig); err != nil {
		return
	}
	return
}
