package person

import (
	"context"
	"ghActionTelegramBot/internal/domain/person"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type storage struct {
	db *mongo.Database
}

func NewStorage(db *mongo.Database) person.Storage {
	return &storage{db: db}
}

func (s *storage) Create(dto person.CreatePersonDto) (*person.Person, error) {
	p := &person.Person{
		ID:          primitive.NewObjectID(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		UserName:    dto.UserName,
		TelegramId:  dto.TelegramId,
		AccessToken: dto.AccessToken,
	}

	_, err := s.db.Collection("users").InsertOne(context.TODO(), p)
	if err != nil {

	}
	return p, nil
}
