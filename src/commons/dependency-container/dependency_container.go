package dependency_container

import (
	"fmt"
	log "go-gallery/src/infrastructure/logger"
	codeGeneratorRepository "go-gallery/src/infrastructure/repository/codeGenerator"
	emailSenderRepository "go-gallery/src/infrastructure/repository/emailSender"
	imageRepository "go-gallery/src/infrastructure/repository/image"
	thumbnailImageRepository "go-gallery/src/infrastructure/repository/image/thumbnailImage"
	userRepository "go-gallery/src/infrastructure/repository/user"
)

type DependencyContainer struct {
	imageRepository          imageRepository.ImageRepository
	userRepository           userRepository.UserRepository
	thumbnailImageRepository thumbnailImageRepository.ThumbnailImageRepository
	codeGeneratorRepository  codeGeneratorRepository.CodeGeneratorRepository
	emailSenderRepository    emailSenderRepository.EmailSenderRepository
}

var dependencyContainer *DependencyContainer

var logger log.Logger

func Instance() *DependencyContainer {
	logger = log.Instance()
	if dependencyContainer == nil {
		return new(DependencyContainer)
	}
	return dependencyContainer
}

func (dp *DependencyContainer) GetUserRepository() userRepository.UserRepository {
	if dp.userRepository != nil {
		return dp.userRepository
	}

	panic("Dependency UserRepository not found.")
}

func (dp *DependencyContainer) SetUserRepository(userDependency userRepository.UserRepository) {
	dp.userRepository = userDependency
	logger.Info(fmt.Sprintf("Dependency UserRepository has been set. Implementation: %T", userDependency))
}

func (dp *DependencyContainer) GetImageRepository() imageRepository.ImageRepository {
	if dp.imageRepository != nil {
		return dp.imageRepository
	}
	panic("Dependency ImageRepository not found.")
}

func (dp *DependencyContainer) SetImageRepository(imageDependency imageRepository.ImageRepository) {
	dp.imageRepository = imageDependency
	logger.Info(fmt.Sprintf("Dependency ImageRepository has been set. Implementation: %T", imageDependency))
}

func (dp *DependencyContainer) GetEmailSenderRepository() emailSenderRepository.EmailSenderRepository {
	if dp.emailSenderRepository != nil {
		return dp.emailSenderRepository
	}
	panic("Dependency EmailSenderRepository not found.")
}

func (dp *DependencyContainer) SetEmailSenderRepository(emailDependency emailSenderRepository.EmailSenderRepository) {
	dp.emailSenderRepository = emailDependency
	logger.Info(fmt.Sprintf("Dependency EmailSenderRepository has been set. Implementation: %T", emailDependency))
}

func (dp *DependencyContainer) SetThumbnailImageRepository(thumbnailImageDependency thumbnailImageRepository.ThumbnailImageRepository) {
	dp.thumbnailImageRepository = thumbnailImageDependency
	logger.Info(fmt.Sprintf("Dependency ThumbnailImageRepository has been set. Implementation: %T", thumbnailImageDependency))
}

func (dp *DependencyContainer) GetThumbnailImageRepository() thumbnailImageRepository.ThumbnailImageRepository {
	if dp.thumbnailImageRepository != nil {
		return dp.thumbnailImageRepository
	}
	panic("Dependency ThumbnailImageRepository not found.")
}

func (dp *DependencyContainer) SetCodeGeneratorRepository(codeGeneratorDependency codeGeneratorRepository.CodeGeneratorRepository) {
	dp.codeGeneratorRepository = codeGeneratorDependency
	logger.Info(fmt.Sprintf("Dependency CodeGeneratorRepository has been set. Implementation: %T", codeGeneratorDependency))
}

func (dp *DependencyContainer) GetCodeGeneratorRepository() codeGeneratorRepository.CodeGeneratorRepository {
	if dp.codeGeneratorRepository != nil {
		return dp.codeGeneratorRepository
	}
	panic("Dependency ThumbnailImageRepository not found.")
}
