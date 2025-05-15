package emailTemplate

type EmailTemplate interface {
	Subject() string
	Body(code string, recipientEmail string) string
}
