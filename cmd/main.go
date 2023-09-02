package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/salesforceanton/pocket-tg-bot/internal/auth_server"
	"github.com/salesforceanton/pocket-tg-bot/internal/config"
	"github.com/salesforceanton/pocket-tg-bot/internal/logger"
	"github.com/salesforceanton/pocket-tg-bot/internal/repository"
	"github.com/salesforceanton/pocket-tg-bot/internal/telegram"
	"github.com/zhashkevych/go-pocket-sdk"
)

func main() {
	// Initilize configs
	cfg, err := config.InitConfig()
	if err != nil {
		logger.LogIssueWithPoint("configs", err)
		return
	}

	// Connect to repo
	db, err := repository.InitBolt(cfg.BoltDBFile)
	if err != nil {
		logger.LogIssueWithPoint("db connect", err)
		return
	}

	// Initialize pocket client
	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKey)
	if err != nil {
		logger.LogIssueWithPoint("pocket client initialization", err)
		return
	}

	// Create auth server instance
	repo := repository.NewTokenStorage(db)
	authServer := auth_server.NewServer(pocketClient, repo, cfg.BotURL)

	// Create new bot instance
	bot, err := telegram.NewBot(cfg, pocketClient, repo)
	if err != nil {
		logger.LogIssueWithPoint("bot creation", err)
		return
	}

	// Run auth server async
	go func() {
		logger.LogInfoWithPoint("auth server", "SERVER RUN SUCCESSFULLY")
		if err := authServer.Run(); err != nil {
			logger.LogIssueWithPoint("auth server running", err)
			return
		}
	}()

	// Start bot async
	go func() {
		logger.LogInfoWithPoint("bot", "BOT RUN SUCCESSFULLY")
		bot.Start()
	}()

	// Graceful shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)
	<-exit

	if err := authServer.Shutdown(context.Background()); err != nil {
		logger.LogIssueWithPoint("shutdown", err)
		return
	}

	if err := db.Close(); err != nil {
		logger.LogIssueWithPoint("shutdown", err)
		return
	}

	logger.LogInfoWithPoint("shutdown", "SERVER SHUTDOWN SUCCESSFULLY")
}
