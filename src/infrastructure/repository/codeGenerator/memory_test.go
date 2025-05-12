package codeGeneratorRepository

import (
	"errors"
	"io"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	USER_EXAMPLE = "sanbricio"
)

var codeGen *CodeGeneratorMemoryRepository

func TestMain(m *testing.M) {
	// Init global instance (como un BeforeAll)
	codeGen = NewCodeGeneratorMemory(make(map[string]string))

	// Run all tests
	code := m.Run()

	os.Exit(code)
}

// Test for GenerateCode
func TestGenerateCode(t *testing.T) {
	// Arrange
	user := USER_EXAMPLE

	// Act
	code, _ := codeGen.GenerateCode(user)

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
	user := USER_EXAMPLE

	// Act
	code, err := codeGen.GenerateCode(user)

	// Assert
	assert.Empty(t, code, "expected empty code, but got a value")
	assert.Error(t, err, "expected error, but got nil")
	assert.ErrorContains(t, err, "simulated random number generator error", "expected simulated random number generator error, but got a different error")
}

// Test for VerifyCode
func TestVerifyCode(t *testing.T) {
	// Arrange
	code, _ := codeGen.GenerateCode(USER_EXAMPLE)

	// Act: Verify the generated code
	isValid := codeGen.VerifyCode(USER_EXAMPLE, code)
	assert.True(t, isValid, "expected code to be valid, but it was invalid")

	// Act: Verify an invalid code
	isValid = codeGen.VerifyCode(USER_EXAMPLE, "invalidCode")
	assert.False(t, isValid, "expected code to be invalid, but it was valid")

	// Act: Verify that an empty code is considered invalid
	isValid = codeGen.VerifyCode(USER_EXAMPLE, "")
	assert.False(t, isValid, "expected code to be invalid, but it was valid")
}

// Test to verify code expiration
func TestVerifyCodeExpired(t *testing.T) {
	// Arrange
	user := USER_EXAMPLE
	code, _ := codeGen.GenerateCode(user)

	// Simulate code expiration (wait more than 5 minutes)
	// Adjust the NowFunc function to return a future time
	NowFunc = func() time.Time {
		return time.Now().Add(6 * time.Minute) // 6 minutes after the current time
	}

	// Act
	isValid := codeGen.VerifyCode(user, code)
	assert.False(t, isValid, "expected code to be expired, but it was valid")

	// Restore the NowFunc function to its original state
	NowFunc = time.Now
}

func TestRemoveCode(t *testing.T) {
	// Arrange
	user := USER_EXAMPLE
	code, _ := codeGen.GenerateCode(user)

	// Act: Remove the code
	codeGen.removeCode(user)

	// Assert: Verify that the code has been removed
	isValid := codeGen.VerifyCode(user, code)
	assert.False(t, isValid, "expected code to be removed, but it was still valid")
}

func TestAutoCleanup(t *testing.T) {
	// Arrange
	user := USER_EXAMPLE
	codeGen := &CodeGeneratorMemoryRepository{
		expirationCode:  100 * time.Millisecond, 
		cleanupInterval: 200 * time.Millisecond, 
	}

	codeGen.GenerateCode(user)
	codeGen.StartAutoCleanup()

	time.Sleep(300 * time.Millisecond)

	// Assert
	mutex.RLock()
	_, exists := codes[user]
	mutex.RUnlock()

	assert.False(t, exists, "expected auto cleanup to remove expired code, but it still exists")
}