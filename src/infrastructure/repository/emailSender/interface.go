package emailSenderRepository

import emailTemplate "go-gallery/src/infrastructure/repository/emailSender/template"

type EmailSenderRepository interface {
	SendEmail(code, email string, template emailTemplate.EmailTemplate) error
}
