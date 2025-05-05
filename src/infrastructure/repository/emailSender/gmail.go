package emailSenderRepository

import (
	"fmt"
	entity "go-gallery/src/domain/entities"
	log "go-gallery/src/infrastructure/logger"
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

func (r *EmailSenderGoMailRepository) SendEmail(code, email string) error {
	logger.Info(fmt.Sprintf("Sending verification email to %s with code: %s", email, code))

	// Mensaje del email
	message := r.buildMessage(code, email)

	// Configuración de autenticación SMTP
	auth := smtp.PlainAuth("", r.emailSender.Username, r.emailSender.Password, r.emailSender.Host)

	// Dirección del servidor SMTP de envío
	serverAddress := fmt.Sprintf("%s:%d", r.emailSender.Host, r.emailSender.Port)

	// Enviar el correo
	err := smtp.SendMail(serverAddress, auth, r.emailSender.From, []string{email}, []byte(message))
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to send email to %s: %s", email, err.Error())
		logger.Error(errorMessage)
		return fmt.Errorf("error sending email: %w", err)
	}

	logger.Info(fmt.Sprintf("Email successfully sent to %s for account deletion", email))
	return nil
}

func (r *EmailSenderGoMailRepository) buildMessage(code, email string) string {
	subject := "⚠️ Verification code to delete your go-gallery account"

	// Cuerpo para el email
	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Account Deletion Confirmation</title>
		</head>
		<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px;">
			<table align="center" width="100%%" bgcolor="#ffffff" style="max-width: 600px; padding: 20px; border-radius: 8px; box-shadow: 0px 0px 10px #cccccc;">
				<tr>
					<td align="center" style="padding-bottom: 20px;">
						<h2 style="color: #333;">🔐 Verification Code</h2>
					</td>
				</tr>
				<tr>
					<td align="center" style="font-size: 16px; color: #555;">
						We received a request to delete your account. To confirm this action, please use the following verification code:
					</td>
				</tr>
				<tr>
					<td align="center" style="padding: 20px;">
						<div style="display: inline-block; background-color: #f8f8f8; padding: 15px 30px; border-radius: 5px; font-size: 24px; font-weight: bold; color: #333; border: 1px solid #ddd;">
							%s
						</div>
					</td>
				</tr>
				<tr>
					<td align="center" style="font-size: 14px; color: #777; padding-top: 20px;">
						⚠️ This code is valid only for the next <strong>5 minutes</strong>. Do not share it with anyone.
					</td>
				</tr>
				<tr>
					<td align="center" style="font-size: 14px; color: #777; padding-top: 10px;">
						If you did not request the deletion of your account, please ignore this message and your account will remain active.
					</td>
				</tr>
				<tr>
					<td align="center" style="padding-top: 30px; font-size: 12px; color: #aaa;">
						Best regards, <br>
						<strong>Support Team</strong><br>
						<a href="mailto:support@yourcompany.com" style="color: #3498db; text-decoration: none;">gogalleryteam@gmail.com</a>
					</td>
				</tr>
			</table>
		</body>
		</html>`, code)

	headers := []string{
		"From: " + r.emailSender.From,
		"To: " + email,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=\"UTF-8\"",
	}

	return strings.Join(headers, "\r\n") + "\r\n\r\n" + htmlBody
}
