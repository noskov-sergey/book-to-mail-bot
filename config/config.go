package config

import (
	"os"
	"strconv"

	"go.uber.org/zap"
)

type Config struct {
	Telegram
	Mail
	Path string
	Size int
}

type Telegram struct {
	Token string
	Host  string
}

type Mail struct {
	Host     string
	Port     string
	From     string
	To       []string
	Password string
}

func MustLoad(logger *zap.Logger) *Config {

	var err error = nil

	var cfg Config

	cfg.Path = os.Getenv("PATH")
	if cfg.Path == "" {
		logger.Fatal("PATH is not set")
	}

	cfg.Size, err = strconv.Atoi(os.Getenv("BATCH_SIZE"))
	if cfg.Size == 0 || err != nil {
		logger.Fatal("BATCH_SIZE is not set")
	}

	cfg.Telegram.Host = os.Getenv("TELEGRAM_HOST")
	if cfg.Telegram.Host == "" {
		logger.Fatal("TELEGRAM_HOST is not set")
	}

	cfg.Telegram.Token = os.Getenv("TELEGRAM_TOKEN")
	if cfg.Telegram.Token == "" {
		logger.Fatal("TELEGRAM_TOKEN is not set")
	}

	cfg.Mail.Port = os.Getenv("MAIL_PORT")
	if cfg.Mail.Port == "" {
		logger.Fatal("MAIL_PORT is not set")
	}

	cfg.Mail.Host = os.Getenv("MAIL_HOST")
	if cfg.Mail.Host == "" {
		logger.Fatal("MAIL_HOST is not set")
	}

	cfg.Mail.Password = os.Getenv("MAIL_PASSWORD")
	if cfg.Mail.Password == "" {
		logger.Fatal("MAIL_PASSWORD is not set")
	}

	cfg.Mail.From = os.Getenv("MAIL_FROM")
	if cfg.Mail.From == "" {
		logger.Fatal("MAIL_FROM is not set")
	}

	cfg.Mail.To = []string{os.Getenv("MAIL_TO")}
	if len(cfg.Mail.To) == 0 {
		logger.Fatal("MAIL_TO is not set")
	}

	logger.Debug("config create successful")

	return &cfg
}
