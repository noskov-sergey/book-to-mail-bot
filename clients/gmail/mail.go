package gmail

import (
	"book-to-mail-bot/lib/e"
	"fmt"
	"github.com/scorredoira/email"
	"net/mail"
	"net/smtp"
)

type Client struct {
	from     string
	password string
	to       string
	host     string
	port     string
}

func New(from string, password string, to string, host string, port string) *Client {
	return &Client{
		from:     from,
		password: password,
		to:       to,
		host:     host,
		port:     port,
	}
}

func (c *Client) SendEmail(file string) error {
	m := email.NewMessage("", "")
	m.From = mail.Address{Name: "From", Address: c.from}
	m.To = c.makeToEmail()

	if err := m.Attach(file); err != nil {
		return e.WrapErr("can't send Email: %w", err)
	}

	m.AddHeader("X-CUSTOMER-id", "xxxxx")

	auth := smtp.PlainAuth("", c.from, c.password, c.host)

	if err := email.Send(c.makeAddress(), auth, m); err != nil {
		return e.WrapErr("can't send Email: %w", err)
	}
	return nil
}

func (c *Client) makeAddress() string {
	return fmt.Sprintf("%s:%s", c.host, c.port)
}

func (c *Client) makeToEmail() []string {
	return []string{c.to}
}
