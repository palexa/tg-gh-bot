package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (b *ghBot) start(msg *tgbotapi.MessageConfig, update *tgbotapi.Update) {
	_, err := b.findOrCreatePerson(update.Message.From)
	if err != nil {
		msg.Text = "smth goes wrong."
	}
}

func (b *ghBot) help(msg *tgbotapi.MessageConfig) {
	msg.Text = "I understand /sayhi and /status and /open and /close."
}

func (b *ghBot) openKeyboard(msg *tgbotapi.MessageConfig) {
	msg.ReplyMarkup = numericKeyboard
}

func (b *ghBot) openInlineKeyboard(msg *tgbotapi.MessageConfig) {
	msg.ReplyMarkup = numericKeyboard2
}

func (b *ghBot) openGitHubAuthKeyboard(msg *tgbotapi.MessageConfig, update *tgbotapi.Update) {
	p, err := b.findOrCreatePerson(update.Message.From)
	if err != nil {
		msg.Text = "Smth goes wrong. Try once more"
	}
	if p.AccessToken != "" {
		msg.Text = "Your github already connected. If you want to unsubscribe, pls type /uns"
		return
	}
	msg.ReplyMarkup = generateGitHubButtonKeyboardMarkup(p.ID.Hex())
}

func (b *ghBot) replyToMessageId(msg *tgbotapi.MessageConfig, messageId int) {
	msg.ReplyToMessageID = messageId
}
