package gmail

import (
	"book-to-mail-bot/lib/e"
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

func (c *Client) SendEmail() error {
	auth := smtp.PlainAuth("", c.from, c.password, c.host)

	err := smtp.SendMail(
		makeAddress(c.host, c.port),
		auth,
		c.from,
		makeToEmail(c.to),
		[]byte{},
	)
	if err != nil {
		return e.WrapErr("can't send Email: %w", err)
	}

	return nil
}

func makeAddress(host string, port string) string {
	return host + ":" + port
}

func makeToEmail(to string) []string {
	return []string{to}
}
