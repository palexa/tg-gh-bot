package config

import (
	"ghActionTelegramBot/middleware"
	"ghActionTelegramBot/models"
	"gopkg.in/yaml.v2"
	"os"
)

var Cfg *models.Config

func LoadConfig() (*models.Config, error) {
	configFile, err := os.Open("config.yml")
	if err != nil {
		middleware.ProcessError(err)
	}
	defer configFile.Close()

	var cfg models.Config
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&cfg)
	if err != nil {
		middleware.ProcessError(err)
	}

	return &cfg, nil

}
