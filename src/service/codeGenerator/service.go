package codeGeneratorService

import (
	"fmt"
	codeGeneratorRepository "go-gallery/src/infrastructure/repository/codeGenerator"
)

type CodeGeneratorService struct {
	repository codeGeneratorRepository.CodeGeneratorRepository
}

func NewCodeGeneratorService(repository codeGeneratorRepository.CodeGeneratorRepository) *CodeGeneratorService {
	return &CodeGeneratorService{
		repository: repository,
	}
}

func (s *CodeGeneratorService) GenerateCode(prefix, key string) (string, error) {
	return s.repository.GenerateCode(fmt.Sprintf("%s-%s", prefix, key))
}

func (s *CodeGeneratorService) VerifyCode(prefix, key, code string) bool {
	return s.repository.VerifyCode(fmt.Sprintf("%s-%s", prefix, key), code)
}
