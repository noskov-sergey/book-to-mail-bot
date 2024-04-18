package clients

type MailClient interface {
	SendEmail(string) error
}
