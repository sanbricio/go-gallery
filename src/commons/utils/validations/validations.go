package validators

import (
	"fmt"
)

const (
	ERROR_MESSAGE   string = "the field '%s' must not be empty"
	UNKNOWN_ERROR   string = "field type '%s' is not supported for required validation"
)

func ValidateNonEmptyStringField(fieldName string, fieldValue any) error {
	switch v := fieldValue.(type) {
	case string:
		if v == "" {
			return fmt.Errorf(ERROR_MESSAGE, fieldName)
		}
	case *string:
		if v == nil || *v == "" {
			return fmt.Errorf(ERROR_MESSAGE, fieldName)
		}
	default:
		return fmt.Errorf(UNKNOWN_ERROR, fieldName)
	}
	return nil
}
