package main

import (
	"log"

	"book-to-mail-bot/clients/gmail"
	tgClient "book-to-mail-bot/clients/telegram"
	"book-to-mail-bot/config"
	event_consumer "book-to-mail-bot/consumer/event-consumer"
	"book-to-mail-bot/events/telegram"
	"book-to-mail-bot/storage/files"
)

func main() {
	cfg := config.MustLoad()

	var eventsProcessor = telegram.New(
		tgClient.New(cfg.Telegram.Host, cfg.Telegram.Token),
		gmail.New(cfg.Mail.From, cfg.Mail.Password, cfg.Mail.To, cfg.Mail.Host, cfg.Mail.Port),
		files.New(cfg.Path),
	)

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, cfg.Size)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
