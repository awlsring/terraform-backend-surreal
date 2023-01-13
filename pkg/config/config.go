package config

import (
	"os"

	"github.com/awlsring/surreal-db-client/surreal"
	"github.com/spf13/viper"
)

type Config struct {
	Port int `mapstructure:"port"`
	Users map[string]string `mapstructure:"users"`
	Surreal surreal.SurrealConfig `mapstructure:"surreal"`
}

func getConfigPath() string {
	path := os.Getenv("CONFIG_PATH")
	if path != "" {
		return path
	}

	return "/config/config.yaml"
}

func LoadConfig() (Config, error) {
	vp := viper.New()
	vp.SetConfigFile(getConfigPath())

	var config Config
	err := vp.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = vp.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}