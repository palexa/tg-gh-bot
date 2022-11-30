package telegram

import (
	"fmt"
	"ghActionTelegramBot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)
var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("1"),
		tgbotapi.NewKeyboardButton("2"),
		tgbotapi.NewKeyboardButton("3"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("4"),
		tgbotapi.NewKeyboardButton("5"),
		tgbotapi.NewKeyboardButton("6"),
	),
)

var numericKeyboard2 = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("dev-vert", "https://a.verticula.xyz"),
		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
	),
)
func InitTelegramBot(config *models.Config) {
	bot, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		fmt.Println("error", err.Error())
		return
	}

	bot.Debug = true
	info,_ := bot.GetWebhookInfo()
	fmt.Println(info)
	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID
		if update.Message != nil {
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "help":
					msg.Text = "I understand /sayhi and /status and /open and /close."
				case "sayhi":
					msg.Text = "Hi :)"
				case "status":
					msg.Text = "I'm ok."
				case "auth":
					msg.Text = "I'm ok."
				case "opens":
					openInlineKeyboard(&msg)
				case "open":
					openKeyboard(&msg)
				case "close":
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				default:
					msg.Text = "I don't know that command"
				}
			} else {
				msg.Text = update.Message.From.FirstName + " " + update.Message.From.LastName
			}

			if _, err := bot.Send(msg); err != nil {
				fmt.Println("error", err.Error())
			}
		}
	}
}

func openKeyboard(msg * tgbotapi.MessageConfig) {
	msg.ReplyMarkup = numericKeyboard
}

func openInlineKeyboard(msg * tgbotapi.MessageConfig) {
	msg.ReplyMarkup = numericKeyboard2
}

func replyToMessageId(msg * tgbotapi.MessageConfig, messageId int) {
	msg.ReplyToMessageID = messageId
}
