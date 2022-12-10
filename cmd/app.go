package main

import (
	"context"
	"fmt"
	task2 "ghActionTelegramBot/internal/adapters/db/task"
	"ghActionTelegramBot/internal/config"
	"ghActionTelegramBot/internal/domain/task"
	"ghActionTelegramBot/pkg/client/mongodb"
	gh_logic2 "ghActionTelegramBot/pkg/gh-logic"
	telegram2 "ghActionTelegramBot/pkg/telegram"
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

	storage := task2.NewStorage(db)

	service := task.NewService(storage)

	t := &task.CreateTaskDto{
		Text: "test str",
	}
	_, err = service.Create(t)
	if err != nil {
		panic(err)
	}
	tasks, _ := service.GetAll()
	fmt.Println(tasks)

	go gh_logic2.RunServer()
	telegram2.InitTelegramBot(config.Cfg)
}

func createTask(task *task.Task) error {
	_, err := collection.InsertOne(context.TODO(), task)
	return err
}
