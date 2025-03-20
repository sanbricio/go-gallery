package emailsender_repository

import (
	"api-upload-photos/src/domain"
	"fmt"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailSenderGoMailRepository struct {
	emailSender domain.EmailSender
}

func NewEmailSenderGoMailRepository(args map[string]string) *EmailSenderGoMailRepository {
	emailSender := domain.EmailSender{}
	// Cargamos de configuración para el sender de emails
	emailSender.Host = args["EMAIL_SENDER_HOST"]
	port, err := strconv.Atoi(args["EMAIL_SENDER_PORT"])
	if err != nil {
		panic(fmt.Sprintf("invalid port value: %v", err))
	}
	emailSender.Port = port
	emailSender.Username = args["EMAIL_SENDER_USERNAME"]
	emailSender.Password = args["EMAIL_SENDER_PASSWORD"]
	emailSender.From = args["EMAIL_SENDER_FROM"]

	return &EmailSenderGoMailRepository{
		emailSender: emailSender,
	}
}

const EmailSenderGoMailRepositoryKey = "EmailSenderGoMailRepository"

// Implementa la interfaz del dominio
func (r *EmailSenderGoMailRepository) SendEmail(code, email string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", r.emailSender.From)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Código de verificación único para eliminar la cuenta")
	m.SetBody("text/plain", "El código de verificación para confirmar la eliminacion de la cuenta es: "+code)

	d := gomail.NewDialer(r.emailSender.Host, r.emailSender.Port, r.emailSender.Username, r.emailSender.Password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("error enviando el email: %w", err)
	}

	return nil
}
