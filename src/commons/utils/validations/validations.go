package validators

import (
	"fmt"
)

const (
	MESSAGE_ERROR string = "el campo '%s' no debe estar vacío"
	UNKNOWN_ERROR string = "tipo de campo '%s' no soportado para validación de requerido"
)

func ValidateNonEmptyStringField(fieldName string, fieldValue interface{}) error {
	switch v := fieldValue.(type) {
	case string:
		if v == "" {
			return fmt.Errorf(MESSAGE_ERROR, fieldName)
		}
	case *string:
		if v == nil || *v == "" {
			return fmt.Errorf(MESSAGE_ERROR, fieldName)
		}
	default:
		return fmt.Errorf(UNKNOWN_ERROR, fieldName)
	}
	return nil
}
