package dependency_dictionary

import (
	"go-gallery/src/infrastructure/logger"
	emailSenderRepository "go-gallery/src/infrastructure/repository/emailSender"
	imageRepository "go-gallery/src/infrastructure/repository/image"
	thumbnailImageRepository "go-gallery/src/infrastructure/repository/image/thumbnailImage"
	userRepository "go-gallery/src/infrastructure/repository/user"
)

func FindLoggerDependency(code string, args map[string]string) logger.Logger {
	//TODO implementar cuando tengamos el de fichero
	return nil
}

func FindImageDependency(code string, args map[string]string) imageRepository.ImageRepository {
	switch code {
	default:
		return imageRepository.NewImageMongoDBRepository(args)
	}
}

func FindThumbnailImageDependency(code string, args map[string]string) thumbnailImageRepository.ThumbnailImageRepository {
	switch code {
	default:
		return thumbnailImageRepository.NewThumbnailImageMongoDBRepository(args)
	}
}

func FindUserDependency(code string, args map[string]string) userRepository.UserRepository {
	switch code {
	case userRepository.UserMongoDBRepositoryKey:
		return userRepository.NewUserMongoDBRepository(args)
	default:
		return userRepository.NewUserPostgreSQLRepository(args)
	}

}

func FindEmailSenderDependency(code string, args map[string]string) emailSenderRepository.EmailSenderRepository {
	switch code {
	default:
		return emailSenderRepository.NewEmailSenderGoMailRepository(args)
	}
}
