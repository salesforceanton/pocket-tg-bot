package telegram

import (
	"context"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/salesforceanton/pocket-tg-bot/internal/logger"
	"github.com/zhashkevych/go-pocket-sdk"
)

const (
	COMMAND_START         = "start"
	REDIRECT_URL_TEMPLATE = "%s?chat_id=%d"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case COMMAND_START:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	chatId := message.Chat.ID
	// Check is user authorized
	accessToken, err := b.repo.GetAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAuthorization(chatId)
	}

	// Saving link into Pocket
	if err = b.putLinkToPocket(message.Text, accessToken); err != nil {
		logger.LogIssueWithPoint("Saving link into Pocket", unableToSaveError)
		return err
	}

	// Reply to message with successful response
	msg := tgbotapi.NewMessage(chatId, b.cfg.LinkSaved)
	msg.ReplyToMessageID = message.MessageID
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	chatId := message.Chat.ID
	// If no Access tocket - init authorization
	_, err := b.repo.GetAccessToken(chatId)
	if err != nil {
		return b.initAuthorization(chatId)
	}

	// If access tocen exist - send message about already authorized
	msg := tgbotapi.NewMessage(chatId, b.cfg.AlreadyAuthorized)
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.cfg.UnknownCommand)
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) putLinkToPocket(link, accessToken string) error {
	// Validate url from message
	err := b.validateURL(link)
	if err != nil {
		logger.LogIssueWithPoint("url validation", invalidUrlError)
		return invalidUrlError
	}

	// Save link from message into Pocket
	return b.pocketClient.Add(context.TODO(), pocket.AddInput{
		URL:         link,
		AccessToken: accessToken,
	})
}

func (b *Bot) validateURL(text string) error {
	_, err := url.ParseRequestURI(text)
	return err
}
