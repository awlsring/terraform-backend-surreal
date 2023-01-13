package config

import (
	"os"

	"github.com/awlsring/surreal-db-client/surreal"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port int `mapstructure:"port"`
	Users map[string]string `mapstructure:"users"`
	Surreal surreal.SurrealConfig `mapstructure:"surreal"`
}

type ConfigMap struct {
	Port int `mapstructure:"port"`
	Surreal surreal.SurrealConfig `mapstructure:"surreal"`
}

type UserMap struct {
	Users map[string]string `mapstructure:"users"`
}

func getConfigPath() string {
	path := os.Getenv("CONFIG_PATH")
	if path != "" {
		return path
	}

	return "/config/config.yaml"
}

func getUsersPath() string {
	path := os.Getenv("USERS_PATH")
	if path != "" {
		return path
	}

	return "/config/users.yaml"
}

func LoadConfigMap() (ConfigMap, error) {
	vp := viper.New()
	vp.SetConfigFile(getConfigPath())

	var config ConfigMap
	err := vp.ReadInConfig()
	if err != nil {
		return ConfigMap{}, err
	}
	err = vp.Unmarshal(&config)
	if err != nil {
		return ConfigMap{}, err
	}
	
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	if user == "" || pass == "" {
		log.Fatalln("DB_USER and DB_PASSWORD must be set in the environment. (or in the config file)")
	}
	config.Surreal.User = user
	config.Surreal.Password = pass
	

	return config, nil
}

func LoadUsersMap() (UserMap, error) {
	vp := viper.New()
	vp.SetConfigFile(getUsersPath())

	var users UserMap
	err := vp.ReadInConfig()
	if err != nil {
		return UserMap{}, err
	}

	err = vp.Unmarshal(&users)
	if err != nil {
		return UserMap{}, err
	}

	return users, nil
}

func LoadConfig() (Config, error) {
	log.Info("Loading config...")


	configMap, err := LoadConfigMap()
	if err != nil {
		log.Fatalln(err)
	}
	log.Infof("Loaded config: %+v", configMap)

	userMap, err := LoadUsersMap()
	if err != nil {
		log.Fatalln(err)
	}
	log.Info("Loaded users")

	config := Config{
		Port: configMap.Port,
		Users: userMap.Users,
		Surreal: configMap.Surreal,
	}

	return config, nil
}