package handler_test

import (
	handler "api-upload-photos/src/infrastructure/controller/handler"
	"testing"
	"time"
)

const (
	EMAIL_EXAMPLE = "valid@example.com"
)

// Test para GenerateCode
func TestGenerateCode(t *testing.T) {
	// Arrange
	email := EMAIL_EXAMPLE

	// Act
	code := handler.GenerateCode(email)

	// Assert
	if code == "" {
		t.Errorf("expected code to be generated, but got empty string")
	}
}

// Test para VerifyCode
func TestVerifyCode(t *testing.T) {
	// Arrange
	code := handler.GenerateCode(EMAIL_EXAMPLE)

	// Act: Verificamos el código generado
	isValid := handler.VerifyCode(EMAIL_EXAMPLE, code)

	// Assert: Verificamos que el código es válido
	if !isValid {
		t.Errorf("expected code to be valid, but it was invalid")
	}

	// Act: Verificamos un código inválido
	isValid = handler.VerifyCode(EMAIL_EXAMPLE, "invalidCode")

	// Assert: El código inválido debe devolver false
	if isValid {
		t.Errorf("expected code to be invalid, but it was valid")
	}

	// Act: Verificamos que un código vacío sea considerado inválido
	isValid = handler.VerifyCode(EMAIL_EXAMPLE, "")

	// Assert: El código vacío debe devolver false
	if isValid {
		t.Errorf("expected code to be invalid, but it was valid")
	}
}

// Test para verificar la expiración del código
func TestVerifyCodeExpired(t *testing.T) {
	// Arrange
	email := EMAIL_EXAMPLE
	code := handler.GenerateCode(email)

	// Simulamos que el código ha expirado (esperamos más de 5 minutos)
	// Ajustamos la función NowFunc para devolver un tiempo futuro
	handler.NowFunc = func() time.Time {
		return time.Now().Add(6 * time.Minute) // 6 minutos después del momento actual
	}

	// Act
	isValid := handler.VerifyCode(email, code)

	// Assert
	if isValid {
		t.Errorf("expected code to be expired, but it was valid")
	}

	// Restaurar la función NowFunc a su estado original
	handler.NowFunc = time.Now
}

func TestRemoveCode(t *testing.T) {
	// Arrange
	email := EMAIL_EXAMPLE
	code := handler.GenerateCode(email)

	// Act: Eliminamos el código
	handler.RemoveCode(email)

	// Assert: Verificamos que el código haya sido eliminado
	isValid := handler.VerifyCode(email, code)
	if isValid {
		t.Errorf("expected code to be removed, but it was still valid")
	}
}
