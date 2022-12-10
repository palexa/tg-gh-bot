package task

import (
	"context"
	"ghActionTelegramBot/internal/domain/task"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type storage struct {
	db *mongo.Database
}

func (s *storage) Create(dto *task.CreateTaskDto) (*task.Task, error) {
	t := &task.Task{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Text:      dto.Text,
		Completed: false,
	}
	_, err := s.db.Collection("tasks").InsertOne(context.TODO(), t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *storage) GetAll() ([]*task.Task, error) {
	// passing bson.D{{}} matches all documents in the collection
	filter := bson.D{{}}
	return s.filterTasks(filter)
}

func NewStorage(db *mongo.Database) task.Storage {
	return &storage{db: db}
}

func (s *storage) filterTasks(filter interface{}) ([]*task.Task, error) {
	var tasks []*task.Task

	cur, err := s.db.Collection("tasks").Find(context.TODO(), filter)
	if err != nil {
		return tasks, err
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var t task.Task
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
