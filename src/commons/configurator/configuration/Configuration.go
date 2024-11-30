package configuration

import (
	"strconv"
	"time"
)

const serviceName = "api-upload-photos"

var configuration *Configuration

type Configuration struct {
	args        map[string]string
	serviceName string
	sesionId    string
	timestamp   time.Time
	port        string
}

func InstanceConfiguration(args map[string]string) *Configuration {
	if configuration == nil {
		timestamp := time.Now()
		miliseconds := timestamp.UnixMilli()
		configuration = &Configuration{}
		configuration.serviceName = serviceName
		configuration.sesionId = serviceName + "-" + strconv.FormatInt(miliseconds, 10)
		configuration.timestamp = timestamp
		configuration.args = args
		configuration.port = args["API_UPLOAD_PHOTOS_PORT"]
		return configuration
	}

	panic("Configuration is already intanced")
}

//TODO Crear getter y setters 
