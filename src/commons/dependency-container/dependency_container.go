package dependency_container

import (
	emailsender_repository "go-gallery/src/infrastructure/repository/emailSenderRepository"
	imageRepository "go-gallery/src/infrastructure/repository/image"
	userRepository "go-gallery/src/infrastructure/repository/user"
	"log"
)

type DependencyContainer struct {
	imageRepository       imageRepository.ImageRepository
	userRepository        userRepository.UserRepository
	emailSenderRepository emailsender_repository.EmailSenderRepository
}

var dependencyContainer *DependencyContainer

func GetIntance() *DependencyContainer {
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
	log.Printf("Dependency UserRepository has been set. Implementation: %T\n", userDependency)
}

func (dp *DependencyContainer) GetImageRepository() imageRepository.ImageRepository {
	if dp.imageRepository != nil {
		return dp.imageRepository
	}
	panic("Dependency ImageRepository not found.")
}

func (dp *DependencyContainer) SetImageRepository(imageDependency imageRepository.ImageRepository) {
	dp.imageRepository = imageDependency
	log.Printf("Dependency ImageRepository has been set. Implementation: %T\n", imageDependency)
}

func (dp *DependencyContainer) GetEmailSenderRepository() emailsender_repository.EmailSenderRepository {
	if dp.emailSenderRepository != nil {
		return dp.emailSenderRepository
	}
	panic("Dependency EmailSenderRepository not found.")
}

func (dp *DependencyContainer) SetEmailSenderRepository(emailDependency emailsender_repository.EmailSenderRepository) {
	dp.emailSenderRepository = emailDependency
	log.Printf("Dependency EmailSenderRepository has been set. Implementation: %T\n", emailDependency)
}
