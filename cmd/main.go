package main

import (
	"github.com/salesforceanton/pocket-tg-bot/internal/config"
	"github.com/salesforceanton/pocket-tg-bot/internal/logger"
	"github.com/salesforceanton/pocket-tg-bot/internal/repository"
	"github.com/salesforceanton/pocket-tg-bot/internal/telegram"
)

func main() {
	// Initilize configs
	cfg, err := config.InitConfig()
	if err != nil {
		logger.LogIssueWithPoint("configs", err)
		return
	}

	// Create and start bot working
	bot, err := telegram.NewBot(cfg)
	if err != nil {
		logger.LogIssueWithPoint("bot creation", err)
		return
	}

	// Connect to repo
	_, err = repository.InitBolt(cfg.BoltDBFile)
	if err != nil {
		logger.LogIssueWithPoint("db connect", err)
		return
	}

	bot.Start()
}
