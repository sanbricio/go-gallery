package entity_test

import (
	entity "go-gallery/src/domain/entities"
	"testing"
)

func TestHashUserPassword(t *testing.T) {
	password := "password123"
	_, err := entity.HashPassword(password)

	if err != nil {
		t.Errorf("TestHashUserPassword: se esperaba que la función HashPassword no devolviera error, pero se obtuvo: %v", err)
		t.FailNow()
	}
}

func TestHashPasswordWithEmptyString(t *testing.T) {
	password := "a very very very very very very very very very very very very long password that exceeds 72 bytes in length"
	_, err := entity.HashPassword(password)

	if err == nil {
		t.Errorf("TestHashPasswordWithEmptyString: se esperaba un error al intentar hashear una contraseña vacía, pero no se obtuvo ninguno.")
		t.FailNow()
	}
}
