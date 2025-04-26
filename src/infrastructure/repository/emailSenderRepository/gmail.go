package emailSenderRepository

import (
	"fmt"
	entity "go-gallery/src/domain/entities"
	"log"
	"net/smtp"
	"strconv"
	"strings"
)

type EmailSenderGoMailRepository struct {
	emailSender entity.EmailSender
}

func NewEmailSenderGoMailRepository(args map[string]string) *EmailSenderGoMailRepository {
	emailSender := entity.EmailSender{}
	// Cargamos de configuración para el sender de emails
	emailSender.Host = args["EMAIL_SENDER_HOST"]
	port, err := strconv.Atoi(args["EMAIL_SENDER_PORT"])
	if err != nil {
		panic(fmt.Sprintf("invalid port value: %v", err))
	}
	emailSender.Port = port
	emailSender.Username = args["EMAIL_SENDER_USERNAME"]
	emailSender.Password = args["EMAIL_SENDER_PASSWORD"]
	emailSender.From = emailSender.Username

	return &EmailSenderGoMailRepository{
		emailSender: emailSender,
	}
}

const EmailSenderGoMailRepositoryKey = "EmailSenderGoMailRepository"

func (r *EmailSenderGoMailRepository) SendEmail(code, email string) error {

	// Mensaje del email
	message := r.buildMessage(code, email)

	// Configuración de autenticación SMTP
	auth := smtp.PlainAuth("", r.emailSender.Username, r.emailSender.Password, r.emailSender.Host)

	// Dirección del servidor SMTP de envío
	serverAddress := fmt.Sprintf("%s:%d", r.emailSender.Host, r.emailSender.Port)

	// Enviar el correo
	err := smtp.SendMail(serverAddress, auth, r.emailSender.From, []string{email}, []byte(message))
	if err != nil {
		log.Printf("ERROR: ❌ No se ha podido enviar el email a %s para la eliminación del usuario", email)
		return fmt.Errorf("error enviando el email: %w", err)
	}

	log.Printf("INFO: ✅ Email enviado correctamente a %s para la eliminación del usuario", email)
	return nil
}

func (r *EmailSenderGoMailRepository) buildMessage(code, email string) string {
	subject := "⚠️ Código de verificación para eliminar tu cuenta de GoGallery"

	// Cuerpo para el email
	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Confirmación de eliminación de cuenta</title>
		</head>
		<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px;">
			<table align="center" width="100%%" bgcolor="#ffffff" style="max-width: 600px; padding: 20px; border-radius: 8px; box-shadow: 0px 0px 10px #cccccc;">
				<tr>
					<td align="center" style="padding-bottom: 20px;">
						<h2 style="color: #333;">🔐 Código de verificación</h2>
					</td>
				</tr>
				<tr>
					<td align="center" style="font-size: 16px; color: #555;">
						Hemos recibido una solicitud para eliminar tu cuenta. Para confirmar esta acción, usa el siguiente código de verificación:
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
						⚠️ Este código es válido solo por los próximos <strong>5 minutos</strong>. No lo compartas con nadie.
					</td>
				</tr>
				<tr>
					<td align="center" style="font-size: 14px; color: #777; padding-top: 10px;">
						Si no solicitaste la eliminación de tu cuenta, ignora este mensaje y tu cuenta permanecerá activa.
					</td>
				</tr>
				<tr>
					<td align="center" style="padding-top: 30px; font-size: 12px; color: #aaa;">
						Atentamente, <br>
						<strong>Equipo de Soporte</strong><br>
						<a href="mailto:support@tuempresa.com" style="color: #3498db; text-decoration: none;">gogalleryteam@gmail.com</a>
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

	// Unimos cabecera y cuerpo
	return strings.Join(headers, "\r\n") + "\r\n\r\n" + htmlBody

}
