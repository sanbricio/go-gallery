package auth

import (
	"testing"
	"time"

	"go-gallery/src/commons/exception"
	log "go-gallery/src/infrastructure/logger"

	jtoken "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func beforeAll() (*JWTTokenManager, string) {
	// init logger
	logger = log.Init(log.NewConsoleLogger())

	// secret config
	secret := "mySecretKey"
	manager := NewJWTTokenManager(secret)

	return manager, secret
}

// Helper function to create a valid token
func createValidToken(manager *JWTTokenManager, username, email string) (string, *exception.ApiException) {
	token, apiErr := manager.CreateToken(username, email)
	if apiErr != nil {
		return "", apiErr
	}
	return token, nil
}

func TestCreateTokenAndValidateSuccess(t *testing.T) {
	manager, _ := beforeAll()

	username := "testuser"
	email := "test@example.com"

	// Create the token and validate
	token, apiErr := createValidToken(manager, username, email)
	assert.Nil(t, apiErr)
	assert.NotEmpty(t, token)

	claims, validateErr := manager.ValidateToken(token)
	assert.Nil(t, validateErr)
	assert.NotNil(t, claims)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, email, claims.Email)

	assert.WithinDuration(t, time.Now(), time.Unix(claims.IssuedAt, 0), 5*time.Second)
	assert.True(t, claims.Expiration > claims.IssuedAt)
}

func TestValidateTokenInvalidSignature(t *testing.T) {
	manager, _ := beforeAll()

	otherManager := NewJWTTokenManager("wrongSecret")
	token, _ := createValidToken(otherManager, "testuser", "test@example.com")

	claims, err := manager.ValidateToken(token)
	assert.Nil(t, claims)
	assert.NotNil(t, err)
	assert.Equal(t, 401, err.Status)
	assert.Contains(t, err.Message, "Invalid token")
}
func TestValidateTokenCorruptedToken(t *testing.T) {
	manager, _ := beforeAll()

	token := "this.is.not.a.valid.token"

	claims, err := manager.ValidateToken(token)
	assert.Nil(t, claims)
	assert.NotNil(t, err)
	assert.Equal(t, 401, err.Status)
	assert.Contains(t, err.Message, "Invalid token")
}

func TestValidateTokenInvalidClaims(t *testing.T) {
	manager, secret := beforeAll()

	customClaims := CustomClaims{
		Foo: "foo",
		Bar: "bar",
		RegisteredClaims: jtoken.RegisteredClaims{
			ExpiresAt: jtoken.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jtoken.NewNumericDate(time.Now()),
		},
	}

	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, customClaims)
	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	// Ahora validamos el token con claims corruptos
	claims, apiErr := manager.ValidateToken(tokenString)

	// Nos aseguramos de que se devuelva el error esperado para claims corruptos
	assert.Nil(t, claims)
	assert.NotNil(t, apiErr)
	assert.Equal(t, 500, apiErr.Status)
	assert.Contains(t, apiErr.Message, "Error in JWT claims")
}

type CustomClaims struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
	jtoken.RegisteredClaims
}
