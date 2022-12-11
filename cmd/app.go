package main

import (
	"context"
	"ghActionTelegramBot/internal/adapters/api/gh"
	"ghActionTelegramBot/internal/adapters/db/person"
	"ghActionTelegramBot/internal/config"
	person2 "ghActionTelegramBot/internal/domain/person"
	"ghActionTelegramBot/pkg/client/mongodb"
	"ghActionTelegramBot/pkg/telegram"
	"log"
)

func main() {
	config.Cfg, _ = config.LoadConfig()

	db, err := mongodb.NewClient(context.Background(),
		config.Cfg.Database.Host,
		config.Cfg.Database.Port,
		config.Cfg.Database.Username,
		config.Cfg.Database.Password,
		config.Cfg.Database.Database,
		config.Cfg.Database.AuthDB,
	)
	if err != nil {
		log.Fatal(err)
	}

	storage := person.NewStorage(db)
	service := person2.NewService(storage)
	handler := gh.NewHandler(service)
	go handler.Run()
	bot, err := telegram.NewTGBot(config.Cfg.Telegram.Token, service)
	if err != nil {
		panic(err)
	}
	bot.Start()
}
