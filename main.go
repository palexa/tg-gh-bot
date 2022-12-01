package main

import (
	"ghActionTelegramBot/config"
	gh_logic "ghActionTelegramBot/gh-logic"
	"ghActionTelegramBot/telegram"
)

func main() {
	config.Cfg, _ = config.LoadConfig()
	go gh_logic.RunServer()
	telegram.InitTelegramBot(config.Cfg)
}
