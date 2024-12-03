package main

import (
	"go.uber.org/zap"

	"book-to-mail-bot/clients/gmail"
	tgClient "book-to-mail-bot/clients/telegram"
	"book-to-mail-bot/config"
	event_consumer "book-to-mail-bot/consumer/event-consumer"
	"book-to-mail-bot/events/telegram"
	"book-to-mail-bot/storage/files"
)

func main() {
	log, _ := zap.NewProduction()

	cfg := config.MustLoad(log)

	var eventsProcessor = telegram.New(
		tgClient.New(cfg.Telegram.Host, cfg.Telegram.Token, log),
		gmail.New(cfg.Mail.From, cfg.Mail.Password, cfg.Mail.To, cfg.Mail.Host, cfg.Mail.Port, log),
		files.New(cfg.Path, log),
		log,
	)

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, cfg.Size, log)

	log.Info("worker is starting")

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", zap.Error(err))
	}
}
