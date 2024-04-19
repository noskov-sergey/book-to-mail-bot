package clients

type MailClient interface {
	SendEmail(f string) error
}
