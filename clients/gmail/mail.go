package gmail

import (
	"fmt"
	"net/mail"
	"net/smtp"

	"github.com/scorredoira/email"
	"go.uber.org/zap"

	"book-to-mail-bot/lib/e"
)

const (
	fromMailTag = "From"
	subjectMail = "Sending book by TG bot"
	bodyMail    = "This book sent from TG API"
)

type Client struct {
	from     string
	password string
	to       []string
	host     string
	port     string

	log *zap.Logger
}

func New(from string, password string, to []string, host string, port string, log *zap.Logger) *Client {
	return &Client{
		from:     from,
		password: password,
		to:       to,
		host:     host,
		port:     port,
		log:      log.Named("email client"),
	}
}

func (c *Client) SendEmail(file string) (err error) {
	defer func() { err = e.WrapIfErr("can't send Email: %w", err) }()

	m := email.NewMessage(subjectMail, bodyMail)
	m.From = mail.Address{Name: fromMailTag, Address: c.from}
	m.To = c.to

	if err := m.Attach(file); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", c.from, c.password, c.host)

	if err := email.Send(c.makeAddress(), auth, m); err != nil {
		return err
	}

	c.log.Info("Book has been sent to mail", zap.String("book name", file), zap.String("email address", c.to[0]))

	return nil
}

func (c *Client) makeAddress() string {
	return fmt.Sprintf("%s:%s", c.host, c.port)
}
