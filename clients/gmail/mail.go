package gmail

import (
	"book-to-mail-bot/lib/e"
	"fmt"
	"github.com/scorredoira/email"
	"log"
	"net/mail"
	"net/smtp"
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
}

func New(from string, password string, to []string, host string, port string) *Client {
	return &Client{
		from:     from,
		password: password,
		to:       to,
		host:     host,
		port:     port,
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

	log.Printf("book '%s' has been sent to mail", file)

	return nil
}

func (c *Client) makeAddress() string {
	return fmt.Sprintf("%s:%s", c.host, c.port)
}
