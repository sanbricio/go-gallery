package codeGeneratorHandler

import (
	"errors"
	"io"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	EMAIL_EXAMPLE = "valid@example.com"
)

// Test for GenerateCode
func TestGenerateCode(t *testing.T) {
	// Arrange
	email := EMAIL_EXAMPLE

	// Act
	code, _ := GenerateCode(email)

	// Assert
	assert.NotEmpty(t, code, "expected code to be generated, but got empty string")
}

func TestGenerateCodeWithError(t *testing.T) {
	// Simulate an error in random code generation
	originalRandFunc := RandFunc
	defer func() { RandFunc = originalRandFunc }() // Restore the original function after the test

	// Inject an error
	RandFunc = func(r io.Reader, max *big.Int) (*big.Int, error) {
		return nil, errors.New("simulated random number generator error")
	}

	// Arrange
	email := EMAIL_EXAMPLE

	// Act
	code, err := GenerateCode(email)

	// Assert
	assert.Empty(t, code, "expected empty code, but got a value")
	assert.Error(t, err, "expected error, but got nil")
	assert.ErrorContains(t, err, "simulated random number generator error", "expected simulated random number generator error, but got a different error")
}

// Test for VerifyCode
func TestVerifyCode(t *testing.T) {
	// Arrange
	code, _ := GenerateCode(EMAIL_EXAMPLE)

	// Act: Verify the generated code
	isValid := VerifyCode(EMAIL_EXAMPLE, code)
	assert.True(t, isValid, "expected code to be valid, but it was invalid")

	// Act: Verify an invalid code
	isValid = VerifyCode(EMAIL_EXAMPLE, "invalidCode")
	assert.False(t, isValid, "expected code to be invalid, but it was valid")

	// Act: Verify that an empty code is considered invalid
	isValid = VerifyCode(EMAIL_EXAMPLE, "")
	assert.False(t, isValid, "expected code to be invalid, but it was valid")
}

// Test to verify code expiration
func TestVerifyCodeExpired(t *testing.T) {
	// Arrange
	email := EMAIL_EXAMPLE
	code, _ := GenerateCode(email)

	// Simulate code expiration (wait more than 5 minutes)
	// Adjust the NowFunc function to return a future time
	NowFunc = func() time.Time {
		return time.Now().Add(6 * time.Minute) // 6 minutes after the current time
	}

	// Act
	isValid := VerifyCode(email, code)
	assert.False(t, isValid, "expected code to be expired, but it was valid")

	// Restore the NowFunc function to its original state
	NowFunc = time.Now
}

func TestRemoveCode(t *testing.T) {
	// Arrange
	email := EMAIL_EXAMPLE
	code, _ := GenerateCode(email)

	// Act: Remove the code
	RemoveCode(email)

	// Assert: Verify that the code has been removed
	isValid := VerifyCode(email, code)
	assert.False(t, isValid, "expected code to be removed, but it was still valid")
}
