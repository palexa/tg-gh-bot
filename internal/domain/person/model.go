package person

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Person struct {
	ID          primitive.ObjectID `bson:"_id"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
	FirstName   string             `bson:"first_name"`
	LastName    string             `bson:"last_name"`
	UserName    string             `bson:"user_name"`
	TelegramId  int64              `bson:"telegram_id"`
	AccessToken string             `bson:"access_token"`
}
