package main

import (
	"context"
	"fmt"
	"ghActionTelegramBot/internal/adapters/db/person"
	"ghActionTelegramBot/internal/config"
	person2 "ghActionTelegramBot/internal/domain/person"
	"ghActionTelegramBot/internal/domain/task"
	"ghActionTelegramBot/pkg/client/mongodb"
	"ghActionTelegramBot/pkg/gh-logic"
	"ghActionTelegramBot/pkg/telegram"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

var collection *mongo.Collection
var ctx = context.TODO()

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

	//taskStorage := task2.NewStorage(db)
	//_ = task.NewService(taskStorage)

	userStorage := person.NewStorage(db)
	service := person2.NewService(userStorage)

	users, err := service.GetAll()
	if err != nil {
		panic(err)
	}
	fmt.Println(users)

	//t := &task.CreateTaskDto{
	//	Text: "test str",
	//}
	//_, err = taskService.Create(t)
	//if err != nil {
	//	panic(err)
	//}
	//tasks, _ := taskService.GetAll()
	//fmt.Println(tasks)

	go gh_logic.RunServer()
	bot, err := telegram.NewTGBot(config.Cfg.Telegram.Token, service)
	if err != nil {
		panic(err)
	}
	bot.Start()
	//telegram.InitTelegramBot(config.Cfg, service)
}

func createTask(task *task.Task) error {
	_, err := collection.InsertOne(context.TODO(), task)
	return err
}
