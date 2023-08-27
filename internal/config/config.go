package config

import (
	"errors"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Config struct {
	TelegramBotToken  string `envconfig:"TELEGRAM_BOT_TOKEN"`
	PocketConsumerKey string `envconfig:"POCKET_CONSUMER_KEY"`
	AuthServerURL     string `envconfig:"AUTH_SERVER_URL"`
	BotURL            string `envconfig:"BOT_URL"`
	BoltDBFile        string `envconfig:"BOLT_DB_FILE"`

	Messages
}

type Messages struct {
	Responses
	Errors
}

type Responses struct {
	Start             string `mapstructure:"start"`
	AlreadyAuthorized string `mapstructure:"already_authorized"`
	UnknownCommand    string `mapstructure:"unknown_command"`
	LinkSaved         string `mapstructure:"link_saved"`
}

type Errors struct {
	Default      string `mapstructure:"default"`
	InvalidURL   string `mapstructure:"invalid_url"`
	UnableToSave string `mapstructure:"unable_to_save"`
}

// Recieve configuration values from env variables
func InitConfig() (*Config, error) {
	var cfg Config

	if err := setUpViper; err != nil {
		return nil, errors.New("Error with config initialization")
	}

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, errors.New("Error with config initialization")
	}

	if err := viper.UnmarshalKey("response", &cfg.Messages.Responses); err != nil {
		return nil, errors.New("Error with config initialization")
	}

	if err := viper.UnmarshalKey("error", &cfg.Messages.Errors); err != nil {
		return nil, errors.New("Error with config initialization")
	}

	return &cfg, nil
}

func setUpViper() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("messages")

	return viper.ReadInConfig()
}
