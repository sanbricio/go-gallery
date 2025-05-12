package emailService

import emailSenderRepository "go-gallery/src/infrastructure/repository/emailSender"

type EmailSenderService struct {
	repository emailSenderRepository.EmailSenderRepository
}

func NewEmailSenderService(repository emailSenderRepository.EmailSenderRepository) *EmailSenderService {
	return &EmailSenderService{
		repository: repository,
	}
}

func (s *EmailSenderService) SendEmail(code, email string) error {
	return s.repository.SendEmail(code, email)
}
