package main

import (
	"context"
	"fmt"
	"ghActionTelegramBot/internal/config"
	"ghActionTelegramBot/pkg/client/mongodb"
	gh_logic2 "ghActionTelegramBot/pkg/gh-logic"
	telegram2 "ghActionTelegramBot/pkg/telegram"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
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
	collection = db.Collection("tasks")
	task := &Task{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Text:      "test str",
		Completed: false,
	}
	err = createTask(task)
	if err != nil {
		panic(err)
	}
	tasks, _ := getAll()
	fmt.Println(tasks)

	go gh_logic2.RunServer()
	telegram2.InitTelegramBot(config.Cfg)
}

type Task struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Text      string             `bson:"text"`
	Completed bool               `bson:"completed"`
}

func createTask(task *Task) error {
	_, err := collection.InsertOne(context.TODO(), task)
	return err
}

func getAll() ([]*Task, error) {
	// passing bson.D{{}} matches all documents in the collection
	filter := bson.D{{}}
	return filterTasks(filter)
}

func filterTasks(filter interface{}) ([]*Task, error) {
	var tasks []*Task

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return tasks, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var t Task
		err := cur.Decode(&t)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, &t)
	}

	if err := cur.Err(); err != nil {
		return tasks, err
	}

	if len(tasks) == 0 {
		return tasks, mongo.ErrNoDocuments
	}

	return tasks, nil
}
