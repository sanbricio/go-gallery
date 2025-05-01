package codeGeneratorHandler

import (
	"errors"
	"io"
	"math/big"
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
	code, _ := GenerateCode(email)

	// Assert
	if code == "" {
		t.Errorf("expected code to be generated, but got empty string")
	}
}

func TestGenerateCodeWithError(t *testing.T) {
	// Simulamos un error en la generación del código aleatorio
	originalRandFunc := RandFunc
	defer func() { RandFunc = originalRandFunc }() // Restaurar la función original después del test

	// Inyectamos un error
	RandFunc = func(r io.Reader, max *big.Int) (*big.Int, error) {
		return nil, errors.New("simulated random number generator error")
	}

	// Arrange
	email := EMAIL_EXAMPLE

	// Act
	code, err := GenerateCode(email)

	// Assert
	if code != "" {
		t.Errorf("expected empty code, but got %s", code)
	}
	if err == nil {
		t.Errorf("expected error, but got nil")
	}
	if err.Error() != "simulated random number generator error" {
		t.Errorf("expected simulated random number generator error, but got %s", err.Error())
	}
}

// Test para VerifyCode
func TestVerifyCode(t *testing.T) {
	// Arrange
	code, _ := GenerateCode(EMAIL_EXAMPLE)

	// Act: Verificamos el código generado
	isValid := VerifyCode(EMAIL_EXAMPLE, code)

	// Assert: Verificamos que el código es válido
	if !isValid {
		t.Errorf("expected code to be valid, but it was invalid")
	}

	// Act: Verificamos un código inválido
	isValid = VerifyCode(EMAIL_EXAMPLE, "invalidCode")

	// Assert: El código inválido debe devolver false
	if isValid {
		t.Errorf("expected code to be invalid, but it was valid")
	}

	// Act: Verificamos que un código vacío sea considerado inválido
	isValid = VerifyCode(EMAIL_EXAMPLE, "")

	// Assert: El código vacío debe devolver false
	if isValid {
		t.Errorf("expected code to be invalid, but it was valid")
	}
}

// Test para verificar la expiración del código
func TestVerifyCodeExpired(t *testing.T) {
	// Arrange
	email := EMAIL_EXAMPLE
	code, _ := GenerateCode(email)

	// Simulamos que el código ha expirado (esperamos más de 5 minutos)
	// Ajustamos la función NowFunc para devolver un tiempo futuro
	NowFunc = func() time.Time {
		return time.Now().Add(6 * time.Minute) // 6 minutos después del momento actual
	}

	// Act
	isValid := VerifyCode(email, code)

	// Assert
	if isValid {
		t.Errorf("expected code to be expired, but it was valid")
	}

	// Restaurar la función NowFunc a su estado original
	NowFunc = time.Now
}

func TestRemoveCode(t *testing.T) {
	// Arrange
	email := EMAIL_EXAMPLE
	code, _ := GenerateCode(email)

	// Act: Eliminamos el código
	RemoveCode(email)

	// Assert: Verificamos que el código haya sido eliminado
	isValid := VerifyCode(email, code)
	if isValid {
		t.Errorf("expected code to be removed, but it was still valid")
	}
}
