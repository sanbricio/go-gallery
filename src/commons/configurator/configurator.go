package configurator

import (
	"fmt"
	"go-gallery/src/commons/configurator/configuration"
	dependency_container "go-gallery/src/commons/dependency-container"
	dependency_dictionary "go-gallery/src/commons/dependency-container/dependency-dictionary"
	"go-gallery/src/infrastructure/logger"

	"os"
	"strings"

	"github.com/joho/godotenv"
)

func LoadConfiguration() (configuration.Configuration, dependency_container.DependencyContainer) {
	rawEnvConfig := loadArgsEnvConfiguration()
	configuration := configuration.Instance(rawEnvConfig)
	logger := buildLogger(configuration)

	logger.Info("Loading configuration...")
	logger.Info(fmt.Sprintf("Session id established: %v", configuration.GetSessionId()))
	logger.Info(fmt.Sprintf("Start date: %v", configuration.GetTimestamp().String()))

	dependencyContainer := buildDependencyContainer(configuration)

	logger.Info("Configuration loaded successfully.")
	return *configuration, *dependencyContainer
}

func loadArgsEnvConfiguration() map[string]string {
	// Cargar variables desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("WARN: Could not load the .env file. Using system environment variables instead.")
	}

	envConfig := make(map[string]string)
	// Recorremos todas las variables de entorno y las establecemos en un map
	for _, rawEnvConfig := range os.Environ() {
		keyValuePair := strings.Split(rawEnvConfig, "=")
		envConfig[keyValuePair[0]] = keyValuePair[1]
	}
	return envConfig
}

func buildLogger(conf *configuration.Configuration) logger.Logger {
	var loggers []logger.Logger

	consoleLogger := logger.NewConsoleLogger()
	loggers = append(loggers, consoleLogger)

	loggerKey := conf.GetArg("LOGGER_TYPE")
	loggerDependency := dependency_dictionary.FindLoggerDependency(loggerKey, conf.GetArgs())
	if loggerDependency != nil {
		loggers = append(loggers, loggerDependency)
	}

	compositeLogger := logger.NewCompositeLogger(loggers...)
	return logger.Init(compositeLogger)
}

func buildDependencyContainer(conf *configuration.Configuration) *dependency_container.DependencyContainer {
	args := conf.GetArgs()
	dp := dependency_container.Instance()

	emailSenderRepositoryKey := conf.GetArg("EMAIL_SENDER_REPOSITORY")
	emailSenderRepositoryDependency := dependency_dictionary.FindEmailSenderDependency(emailSenderRepositoryKey, args)
	dp.SetEmailSenderRepository(emailSenderRepositoryDependency)

	userRepositoryKey := conf.GetArg("USER_REPOSITORY")
	userRepositoryDependency := dependency_dictionary.FindUserDependency(userRepositoryKey, args)
	dp.SetUserRepository(userRepositoryDependency)

	imageRepositoryKey := conf.GetArg("IMAGE_REPOSITORY")
	imageRepositoryDependency := dependency_dictionary.FindImageDependency(imageRepositoryKey, args)
	dp.SetImageRepository(imageRepositoryDependency)

	thumbnailImageRepositoryKey := conf.GetArg("THUMBNAIL_IMAGE_REPOSITORY")
	thumbnailImageRepositoryDependency := dependency_dictionary.FindThumbnailImageDependency(thumbnailImageRepositoryKey, args)
	dp.SetThumbnailImageRepository(thumbnailImageRepositoryDependency)

	return dp
}
