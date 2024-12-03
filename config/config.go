package config

import (
	"log"
	"os"
	"strconv"
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

func MustLoad() *Config {

	var err error = nil

	var cfg Config

	cfg.Path = os.Getenv("PATH")
	if cfg.Path == "" {
		log.Fatal("PATH is not set")
	}

	cfg.Size, err = strconv.Atoi(os.Getenv("BATCH_SIZE"))
	if cfg.Size == 0 || err != nil {
		log.Fatal("BATCH_SIZE is not set")
	}

	cfg.Telegram.Host = os.Getenv("TELEGRAM_HOST")
	if cfg.Telegram.Host == "" {
		log.Fatal("TELEGRAM_HOST is not set")
	}

	cfg.Telegram.Token = os.Getenv("TELEGRAM_TOKEN")
	if cfg.Telegram.Token == "" {
		log.Fatal("TELEGRAM_TOKEN is not set")
	}

	cfg.Mail.Port = os.Getenv("MAIL_PORT")
	if cfg.Mail.Port == "" {
		log.Fatal("MAIL_PORT is not set")
	}

	cfg.Mail.Host = os.Getenv("MAIL_HOST")
	if cfg.Mail.Host == "" {
		log.Fatal("MAIL_HOST is not set")
	}

	cfg.Mail.Password = os.Getenv("MAIL_PASSWORD")
	if cfg.Mail.Password == "" {
		log.Fatal("MAIL_PASSWORD is not set")
	}

	cfg.Mail.From = os.Getenv("MAIL_FROM")
	if cfg.Mail.From == "" {
		log.Fatal("MAIL_FROM is not set")
	}

	cfg.Mail.To = []string{os.Getenv("MAIL_TO")}
	if len(cfg.Mail.To) == 0 {
		log.Fatal("MAIL_TO is not set")
	}

	return &cfg
}
