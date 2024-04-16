package clients

type MailClient interface {
	SendEmail() error
}
