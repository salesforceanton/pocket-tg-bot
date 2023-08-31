package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/salesforceanton/pocket-tg-bot/internal/logger"
)

func (b *Bot) initAuthorization(chatId int64) error {
	redirectUrl := b.generateRedirectUrl(chatId)

	// Generate auth link
	authLink, err := b.generateAuthLink(chatId, redirectUrl)
	if err != nil {
		logger.LogIssueWithPoint("error with saving request token", err)
		return err
	}

	// Send message to target chat
	msgText := fmt.Sprintf(b.cfg.Messages.Start, authLink)
	msg := tgbotapi.NewMessage(chatId, msgText)
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) generateAuthLink(chatId int64, redirectUrl string) (string, error) {
	// Get request token from pocket API
	requestToken, err := b.pocketClient.GetRequestToken(context.TODO(), redirectUrl)
	if err != nil {
		logger.LogIssueWithPoint("error with getting request token", err)
		return "", err
	}

	// Store request token in repository
	err = b.repo.SaveRequestToken(chatId, requestToken)
	if err != nil {
		logger.LogIssueWithPoint("error with saving request token", err)
		return "", err
	}

	// Generate auth link
	authLink, err := b.pocketClient.GetAuthorizationURL(requestToken, redirectUrl)
	if err != nil {
		logger.LogIssueWithPoint("error with saving request token", err)
		return "", err
	}

	return authLink, nil
}

func (b *Bot) generateRedirectUrl(chatId int64) string {
	return fmt.Sprintf(REDIRECT_URL_TEMPLATE, b.cfg.AuthServerURL, chatId)
}
