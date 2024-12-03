package event_consumer

import (
	"time"

	"go.uber.org/zap"

	"book-to-mail-bot/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int

	log *zap.Logger
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int, log *zap.Logger) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
		log:       log.Named("worker"),
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			c.log.Error("fetch:", zap.Error(err))

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err = c.handleEvents(gotEvents); err != nil {
			c.log.Error("handle events:", zap.Error(err))

			continue
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		c.log.Info("got new event", zap.String("event text", event.Text))

		if err := c.processor.Process(event); err != nil {
			c.log.Error("can't handle event:", zap.Error(err))

			continue
		}
	}

	return nil
}
