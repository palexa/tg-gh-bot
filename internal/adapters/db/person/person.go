package person

import (
	"context"
	"fmt"
	"ghActionTelegramBot/internal/domain/person"
	"go.mongodb.org/mongo-driver/bson"
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

func (s *storage) FindByTelegramId(telegramId int) (*person.Person, error) {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"telegram_id", bson.D{{"$eq", 25}}},
				},
			},
		},
	}
	cursor, err := s.db.Collection("users").Find(context.TODO(), filter)
	if err != nil {

	}
	var result bson.M
	// check for errors in the finding
	if err = cursor.Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println(result)
	return nil, nil
}

func (s *storage) Create(dto *person.CreatePersonDto) (*person.Person, error) {
	p := &person.Person{
		ID:         primitive.NewObjectID(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		FirstName:  dto.FirstName,
		LastName:   dto.LastName,
		UserName:   dto.UserName,
		TelegramId: dto.TelegramId,
	}

	_, err := s.db.Collection("users").InsertOne(context.TODO(), p)
	if err != nil {

	}
	return p, nil
}

func (s *storage) Update(dto *person.UpdatePersonDto) (*person.Person, error) {
	p := &person.Person{
		UpdatedAt:   time.Now(),
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		UserName:    dto.UserName,
		TelegramId:  dto.TelegramId,
		AccessToken: dto.AccessToken,
	}
	s.db.Collection("users").UpdateByID(context.TODO(), dto.ID, p)
	return nil, nil
}

func (s *storage) GetAll() ([]*person.Person, error) {
	filter := bson.D{{}}
	return s.filterUsers(filter)
}

func (s *storage) filterUsers(filter interface{}) ([]*person.Person, error) {
	var users []*person.Person

	cur, err := s.db.Collection("users").Find(context.TODO(), filter)
	if err != nil {
		return users, err
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var u person.Person
		err := cur.Decode(&u)
		if err != nil {
			return users, err
		}

		users = append(users, &u)
	}

	if err := cur.Err(); err != nil {
		return users, err
	}

	if len(users) == 0 {
		return users, mongo.ErrNoDocuments
	}

	return users, nil
}

func (s *storage) filterOne(filter interface{}) (*person.Person, error) {
	var user person.Person
	result := s.db.Collection("users").FindOne(context.TODO(), filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, result.Err()
		}
		panic(result.Err())
	}
	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *storage) FindOrCreate(dto *person.CreatePersonDto) (*person.Person, error) {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"telegram_id", bson.D{{"$eq", dto.TelegramId}}},
				},
			},
		},
	}
	u, err := s.filterOne(filter)
	if err != nil {
		panic(err)
	}
	return u, err
}
