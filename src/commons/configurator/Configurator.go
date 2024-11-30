package configurator

import (
	"api-upload-photos/src/commons/configurator/configuration"
	dependency_container "api-upload-photos/src/commons/dependency-container"
	"log"
	"os"
	"strings"
)

func LoadConfiguration() (configuration.Configuration, dependency_container.DependencyContainer) {
	envConfig := loadArgsEnvConfiguration()
	configuration := configuration.InstanceConfiguration(envConfig)
	//TODO BUILDDEPENDENCYCONTAINER
	//TODO Logs para indicar que se esta iniciando la configuracion 
}

func loadArgsEnvConfiguration() map[string]string {
	envConfig := make(map[string]string)
	// Recorremos todas las variables de entorno y las establecemos en un map
	for _, rawEnvConfig := range os.Environ() {
		keyValuePair := strings.Split(rawEnvConfig, "=")
		envConfig[keyValuePair[0]] = keyValuePair[1]
	}
	return envConfig
}

