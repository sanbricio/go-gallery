package emailService

import (
	emailSenderRepository "go-gallery/src/infrastructure/repository/emailSender"
	emailTemplate "go-gallery/src/infrastructure/repository/emailSender/template"
)

type EmailSenderService struct {
	repository emailSenderRepository.EmailSenderRepository
}

func NewEmailSenderService(repository emailSenderRepository.EmailSenderRepository) *EmailSenderService {
	return &EmailSenderService{
		repository: repository,
	}
}

func (s *EmailSenderService) SendEmail(code, email string, template emailTemplate.EmailTemplate) error {
	return s.repository.SendEmail(code, email, template)
}
