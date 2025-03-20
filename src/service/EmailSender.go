package service

import emailsender_repository "api-upload-photos/src/infrastructure/repository/emailSender"

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
