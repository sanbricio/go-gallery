package exception

import "fmt"

type BuilderException struct {
	Field   string
	Message string
}

func NewBuilderException(fieldName, message string) *BuilderException {
	return &BuilderException{
		Field:   fieldName,
		Message: message,
	}
}

func (e *BuilderException) Error() string {
	return fmt.Sprintf("Error en el campo '%s': %s", e.Field, e.Message)
}
