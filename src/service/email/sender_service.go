package emailService

import emailsender_repository "go-gallery/src/infrastructure/repository/emailSenderRepository"

type EmailSenderService struct {
	repository emailsender_repository.EmailSenderRepository
}

func NewEmailSenderService(repository emailsender_repository.EmailSenderRepository) *EmailSenderService {
	return &EmailSenderService{
		repository: repository,
	}
}

func (s *EmailSenderService) SendEmail(code, email string) error {
	return s.repository.SendEmail(code, email)
}
