package telegram

import (
	"log"
	"strings"

	"book-to-mail-bot/lib/e"
	"book-to-mail-bot/storage"
)

const (
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string, fileID string, fileName string) error {
	if fileID != "" {
		log.Printf("got new book '%s' from '%s'", fileName, username)

		return p.sendBook(chatID, fileID, fileName)
	}

	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	switch text {
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) sendBook(chatID int, id string, name string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: send book", err) }()

	fileLink, err := p.tg.GetFileLink(id)
	if err != nil {
		return err
	}

	file, err := p.tg.DownloadFile(fileLink.Path)
	if err != nil {
		return err
	}

	book := &storage.Book{
		Name: name,
	}

	fPath, err := p.storage.Save(book, file)
	if err != nil {
		return err
	}

	if err := p.mail.SendEmail(fPath); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}
