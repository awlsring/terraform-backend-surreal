package config

import (
	"github.com/awlsring/surreal-db-client/surreal"
	"github.com/spf13/viper"
)

type Config struct {
	Port int `mapstructure:"port"`
	Users map[string]string `mapstructure:"users"`
	Surreal surreal.SurrealConfig `mapstructure:"surreal"`
}

func LoadConfig() (Config, error) {
	vp := viper.New()

	var config Config
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.AddConfigPath(".")
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