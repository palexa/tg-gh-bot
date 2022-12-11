package telegram

import (
	"ghActionTelegramBot/internal/domain/person"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *ghBot) findOrCreatePerson(user *tgbotapi.User) (*person.Person, error) {
	u, err := b.service.FindOrCreate(&person.CreatePersonDto{
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		UserName:   user.UserName,
		TelegramId: user.ID,
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}
