package telegram

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	invalidUrlError   = errors.New("url is invalid")
	unableToSaveError = errors.New("unable to save link to Pocket")
)

func (b *Bot) handleError(chatID int64, err error) {
	var messageText string

	switch err {
	case invalidUrlError:
		messageText = b.config.Messages.Errors.InvalidURL
	case unableToSaveError:
		messageText = b.config.Messages.Errors.UnableToSave
	default:
		messageText = b.config.Messages.Errors.Default
	}

	msg := tgbotapi.NewMessage(chatID, messageText)
	b.bot.Send(msg)
}
