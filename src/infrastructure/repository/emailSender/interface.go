package emailSenderRepository

type EmailSenderRepository interface {
	SendEmail(code, email string) error
}
