package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/salesforceanton/pocket-tg-bot/internal/config"
)

type Bot struct {
	config *config.Config
	bot    *tgbotapi.BotAPI
}

func NewBot(cfg *config.Config) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		return nil, err
	}

	bot.Debug = true
	return &Bot{
		config: cfg,
		bot:    bot,
	}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Run cycle in goroutine to ping tg API for updates
	updates := b.bot.GetUpdatesChan(u)

	// Handle Updates
	for update := range updates {
		// If we got a message
		if update.Message == nil {
			continue
		}

		chatId := update.Message.Chat.ID
		// Handle Command
		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(chatId, err)
			}
		}

		// Handle regular messages
		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}
	}

}
