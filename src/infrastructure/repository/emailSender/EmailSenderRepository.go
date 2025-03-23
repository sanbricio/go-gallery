package emailsender_repository

type EmailSenderRepository interface {
	SendEmail(code, email string) error
}
