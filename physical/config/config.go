package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func New(fileName string) *viper.Viper {
	config := viper.New()

	config.SetConfigName(fileName)

	config.AddConfigPath(".")
	config.AddConfigPath("$HOME")

	err := config.ReadInConfig()
	if err != nil {
		log.Fatal("Error while parsing configuration file", err)
	}

	return config
}

func GetConfigFileName() string {
	env := os.Getenv("ENV")

	if env != "" {
		return "arcbox-" + env
	}

	return "arcbox"
}
