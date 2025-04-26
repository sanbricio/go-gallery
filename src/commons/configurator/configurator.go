package configurator

import (
	"go-gallery/src/commons/configurator/configuration"
	dependency_container "go-gallery/src/commons/dependency-container"
	dependency_dictionary "go-gallery/src/commons/dependency-container/dependency-dictionary"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func LoadConfiguration() (configuration.Configuration, dependency_container.DependencyContainer) {
	rawEnvConfig := loadArgsEnvConfiguration()
	configuration := configuration.Instance(rawEnvConfig)

	log.Println("INFO: Loading configuration...")
	log.Println("INFO: Session id established:", configuration.GetSessionId())
	log.Println("INFO: Start date:", configuration.GetTimestamp().String())

	dependencyContainer := buildDependencyContainer(configuration)

	log.Println("INFO: Configuration loaded successfully.")
	return *configuration, *dependencyContainer
}

func loadArgsEnvConfiguration() map[string]string {
	// Cargar variables desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("WARN: Could not load the .env file. Using system environment variables instead.")
	}

	envConfig := make(map[string]string)
	// Recorremos todas las variables de entorno y las establecemos en un map
	for _, rawEnvConfig := range os.Environ() {
		keyValuePair := strings.Split(rawEnvConfig, "=")
		envConfig[keyValuePair[0]] = keyValuePair[1]
	}
	return envConfig
}

func buildDependencyContainer(conf *configuration.Configuration) *dependency_container.DependencyContainer {
	args := conf.GetArgs()
	dependencyContainer := dependency_container.GetIntance()

	emailSenderRepositoryKey := conf.GetArg("EMAIL_SENDER_REPOSITORY")
	emailSenderRepositoryDependency := dependency_dictionary.FindEmailSenderDependency(emailSenderRepositoryKey, args)
	dependencyContainer.SetEmailSenderRepository(emailSenderRepositoryDependency)

	userRepositoryKey := conf.GetArg("USER_REPOSITORY")
	userRepositoryDependency := dependency_dictionary.FindUserDependency(userRepositoryKey, args)
	dependencyContainer.SetUserRepository(userRepositoryDependency)

	imageRepositoryKey := conf.GetArg("IMAGE_REPOSITORY")
	imageRepositoryDependency := dependency_dictionary.FindImageDependency(imageRepositoryKey, args)
	dependencyContainer.SetImageRepository(imageRepositoryDependency)

	return dependencyContainer
}
