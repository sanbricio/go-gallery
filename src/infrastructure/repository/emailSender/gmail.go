package emailSenderRepository

import (
	"fmt"
	entity "go-gallery/src/domain/entities"
	log "go-gallery/src/infrastructure/logger"
	emailTemplate "go-gallery/src/infrastructure/repository/emailSender/template"
	"net/smtp"
	"strconv"
	"strings"
)

const EmailSenderGoMailRepositoryKey = "EmailSenderGoMailRepository"

var logger log.Logger

type EmailSenderGoMailRepository struct {
	emailSender entity.EmailSender
}

func NewEmailSenderGoMailRepository(args map[string]string) *EmailSenderGoMailRepository {
	emailSender := entity.EmailSender{}
	logger = log.Instance()

	logger.Info("Initializing EmailSenderGoMailRepository...")

	emailSender.Host = args["EMAIL_SENDER_HOST"]
	port, err := strconv.Atoi(args["EMAIL_SENDER_PORT"])
	if err != nil {
		panicMessage := fmt.Sprintf("Invalid value for port: %v", err)
		logger.Panic(panicMessage)
		panic(panicMessage)
	}

	emailSender.Port = port
	emailSender.Username = args["EMAIL_SENDER_USERNAME"]
	emailSender.Password = args["EMAIL_SENDER_PASSWORD"]
	emailSender.From = emailSender.Username

	logger.Info(fmt.Sprintf("EmailSender successfully configured with host: %s and port: %d", emailSender.Host, emailSender.Port))

	return &EmailSenderGoMailRepository{
		emailSender: emailSender,
	}
}

func (r *EmailSenderGoMailRepository) SendEmail(code, email string, template emailTemplate.EmailTemplate) error {
	logger.Info(fmt.Sprintf("Sending verification email to %s with code: %s", email, code))

	headers := []string{
		"From: " + r.emailSender.From,
		"To: " + email,
		"Subject: " + template.Subject(),
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=\"UTF-8\"",
	}

	// email message
	message := strings.Join(headers, "\r\n") + "\r\n\r\n" + template.Body(code, email)

	// SMTP auth
	auth := smtp.PlainAuth("", r.emailSender.Username, r.emailSender.Password, r.emailSender.Host)

	// SMTP server address
	serverAddress := fmt.Sprintf("%s:%d", r.emailSender.Host, r.emailSender.Port)

	err := smtp.SendMail(serverAddress, auth, r.emailSender.From, []string{email}, []byte(message))
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to send email to %s: %s", email, err.Error())
		logger.Error(errorMessage)
		return fmt.Errorf("error sending email: %w", err)
	}

	logger.Info(fmt.Sprintf("Email successfully sent to %s for account deletion", email))
	return nil
}