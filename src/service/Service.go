package service

import (
	"api-upload-photos/src/domain"
	infrastructure "api-upload-photos/src/infrastructure/repository"
	"mime/multipart"
)

type Service struct {
	repository infrastructure.IRepository
}

func NewService(repository infrastructure.IRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Insert(fileInput *multipart.FileHeader) (*domain.Response, error) {
	return s.repository.Insert(fileInput)
}

func (s *Service) Find(id string) (*domain.DTOImage, error){
	return s.repository.Find(id)
}
