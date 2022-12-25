package config

import (
	"ghActionTelegramBot/pkg/logging"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	BaseUrl  string `yaml:"base_url"`
	Telegram struct {
		Token string `yaml:"token"`
	} `yaml:"telegram"`
	GitHub struct {
		ClientId     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
	} `yaml:"git_hub"`
	Database struct {
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		Host       string `yaml:"host"`
		Port       string `yaml:"port"`
		Collection string `yaml:"collection"`
		Database   string `yaml:"database"`
		AuthDB     string `yaml:"auth_db"`
	} `yaml:"database"`
}

var Cfg *Config

func LoadConfig() (*Config, error) {
	configFile, err := os.Open("config.yml")
	if err != nil {
		logging.ProcessError(err)
	}
	defer configFile.Close()

	var cfg Config
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&cfg)
	if err != nil {
		logging.ProcessError(err)
	}

	return &cfg, nil
}
