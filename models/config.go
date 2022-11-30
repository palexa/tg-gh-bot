package models

type Config struct {
	Telegram struct {
		Token string `yaml:"token"`
	} `yaml:"telegram"`
	GitHub struct {
		ClientId string `yaml:"clientId"`
		ClientSecret string `yaml:"clientSecret"`
	} `yaml:"gitHub"`
	Database struct {
		Username string `yaml:"user"`
		Password string `yaml:"pass"`
	} `yaml:"database"`
}
