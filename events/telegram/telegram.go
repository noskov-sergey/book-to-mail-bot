package telegram

import (
	"book-to-mail-bot/clients"
	"book-to-mail-bot/clients/telegram"
	"book-to-mail-bot/events"
	"book-to-mail-bot/lib/e"
	"errors"
)

type Processor struct {
	tg     telegram.Client
	mail   clients.MailClient
	offset int
}

type Meta struct {
	ChatID   int
	Username string
	FileID   string
	FilName  string
}

var ErrUnknownEventType = errors.New("unknown event type")
var ErrUnknownMetaType = errors.New("unknown meta type")

func New(tg telegram.Client, mail clients.MailClient) Processor {
	return Processor{
		tg:   tg,
		mail: mail,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.WrapErr("can't get events: %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Document, events.Message:
		return p.processMessage(event)
	default:
		return e.WrapErr("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) (err error) {
	defer func() { err = e.WrapIfErr("can't process message: %w", err) }()

	meta, err := meta(event)

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username, meta.FileID, meta.FilName); err != nil {
		return err
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.WrapErr("can't get meta", ErrUnknownMetaType)
	}
	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: fetchType(upd),
		Text: fetchText(upd),
	}

	if updType == events.Document {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
			FileID:   upd.Message.Document.ID,
			FilName:  upd.Message.Document.Name,
		}
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	if upd.Message.Document != nil {
		return events.Document
	}

	return events.Message
}
