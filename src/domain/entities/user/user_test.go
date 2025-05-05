package userEntity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashUserPassword(t *testing.T) {
	password := "password123"
	_, err := HashPassword(password)

	assert.NoError(t, err, "Se esperaba que la función HashPassword no devolviera error")
}

func TestHashPasswordWithEmptyString(t *testing.T) {
	password := "a very very very very very very very very very very very very long password that exceeds 72 bytes in length"
	_, err := HashPassword(password)

	assert.Error(t, err, "Se esperaba un error al intentar hashear una contraseña vacía")
}
