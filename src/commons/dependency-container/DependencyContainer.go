package dependency_container

import (
	image_repository "api-upload-photos/src/infrastructure/repository/image"
	user_repository "api-upload-photos/src/infrastructure/repository/user"
	"log"
)

type DependencyContainer struct {
	imageRepository image_repository.ImageRepository
	userRepository  user_repository.UserRepository
}

var dependencyContainer *DependencyContainer

func GetIntance() *DependencyContainer {
	if dependencyContainer == nil {
		return new(DependencyContainer)
	}
	return dependencyContainer
}

func (dp *DependencyContainer) GetUserRepository() user_repository.UserRepository {
	if dp.userRepository != nil {
		return dp.userRepository
	}

	panic("Dependency UserRepository not found.")
}

func (dp *DependencyContainer) SetUserRepository(userDependency user_repository.UserRepository) {
	dp.userRepository = userDependency
	log.Printf("Dependency UserRepository has been set. Implementation: %T\n", userDependency)
}

func (dp *DependencyContainer) GetImageRepository() image_repository.ImageRepository {
	if dp.imageRepository != nil {
		return dp.imageRepository
	}
	panic("Dependency ImageRepository not found.")
}

func (dp *DependencyContainer) SetImageRepository(imageDependency image_repository.ImageRepository) {
	dp.imageRepository = imageDependency
	log.Printf("Dependency ImageRepository has been set. Implementation: %T\n", imageDependency)
}
