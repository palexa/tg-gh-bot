package telegram

import (
	"fmt"
	"ghActionTelegramBot/internal/domain/person"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TGBot interface {
	Start()
}

type ghBot struct {
	bot     *tgbotapi.BotAPI
	service person.Service
}

func NewTGBot(tgToken string, service person.Service) (TGBot, error) {
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		fmt.Println("error", err.Error())
		return nil, err
	}

	bot.Debug = true
	info, _ := bot.GetWebhookInfo()
	fmt.Println(info)
	return &ghBot{service: service, bot: bot}, nil
}

func (b *ghBot) Start() {
	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30
	updates := b.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		if update.Message != nil {
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					b.start(&msg, &update)
				case "help":
					b.help(&msg)
				case "sayhi":
					msg.Text = "Hi :)"
				case "status":
					msg.Text = "I'm ok."
				case "auth":
					msg.Text = "I'm ok."
				case "opens":
					b.openInlineKeyboard(&msg)
				case "open":
					b.openKeyboard(&msg)
				case "close":
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				default:
					msg.Text = "I don't know that command"
				}
			} else {
				msg.Text = update.Message.From.FirstName + " " + update.Message.From.LastName
			}

			if _, err := b.bot.Send(msg); err != nil {
				fmt.Println("error", err.Error())
			}
		}
	}
}
