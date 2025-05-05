package userMiddleware

import (
	"go-gallery/src/infrastructure/auth"
	log "go-gallery/src/infrastructure/logger"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const (
	EMAIL_EXAMPLE    = "test@example.com"
	USERNAME_EXAMPLE = "testuser"
	SECRET           = "testsecret"
	TEST_ROUTE       = "/test"
)

func beforeAll() {
	logger = log.Init(log.NewConsoleLogger())
}

func setupApp(secret string, handler fiber.Handler) (*fiber.App, *JWTMiddleware, *auth.JWTTokenManager) {
	tokenManager := auth.NewJWTTokenManager(secret)
	middleware := NewJWTMiddleware(tokenManager)
	app := fiber.New()
	app.Post(TEST_ROUTE, middleware.Handler(), handler)
	return app, middleware, tokenManager
}

func sendTestRequest(app *fiber.App, token string) (*http.Response, error) {
	req := httptest.NewRequest("POST", TEST_ROUTE, nil)
	if token != "" {
		req.Header.Set("Cookie", COOKIE_NAME+"="+token)
	}
	return app.Test(req)
}

func TestHandlerWithValidToken(t *testing.T) {
	beforeAll()
	app, middleware, _ := setupApp(SECRET, func(c *fiber.Ctx) error {
		return c.SendString("Authenticated")
	})

	validToken, apiErr := middleware.tokenManager.CreateToken(USERNAME_EXAMPLE, EMAIL_EXAMPLE)
	assert.Nil(t, apiErr)

	resp, err := sendTestRequest(app, validToken)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHandlerWithNoToken(t *testing.T) {
	beforeAll()
	app, _, _ := setupApp(SECRET, func(c *fiber.Ctx) error {
		return c.SendString("Authenticated")
	})

	resp, err := sendTestRequest(app, "")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestHandlerWithInvalidToken(t *testing.T) {
	beforeAll()
	app, _, _ := setupApp(SECRET, func(c *fiber.Ctx) error {
		return c.SendString("Authenticated")
	})

	resp, err := sendTestRequest(app, "invalidtokenformat")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestHandlerWithExpiringTokenAndShouldRenew(t *testing.T) {
	beforeAll()
	app, _, _ := setupApp(SECRET, func(c *fiber.Ctx) error {
		newCookie := c.Cookies(COOKIE_NAME)
		assert.NotEmpty(t, newCookie, "Expected cookie to be renewed, but got empty")
		return c.SendString("Authenticated")
	})

	// Creamos token que expira en 5 minutos para simular que está a punto de expirar
	claims := jtoken.MapClaims{
		"username": USERNAME_EXAMPLE,
		"email":    EMAIL_EXAMPLE,
		"exp":      time.Now().Add(5 * time.Minute).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
	expiringToken, err := token.SignedString([]byte(SECRET))
	assert.NoError(t, err)

	resp, err := sendTestRequest(app, expiringToken)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDeleteAuthCookie(t *testing.T) {
	beforeAll()
	app, middleware, tokenManager := setupApp(SECRET, func(c *fiber.Ctx) error {
		return c.SendString("Authenticated")
	})

	validToken, apiErr := middleware.tokenManager.CreateToken(USERNAME_EXAMPLE, EMAIL_EXAMPLE)
	assert.Nil(t, apiErr)

	_, err := sendTestRequest(app, validToken)
	assert.NoError(t, err)

	tokenManagerMiddleware := NewJWTMiddleware(tokenManager)
	app.Post("/delete", func(c *fiber.Ctx) error {
		tokenManagerMiddleware.DeleteAuthCookie(c)
		return c.SendString("Cookie deleted")
	})

	req := httptest.NewRequest("POST", "/delete", nil)
	req.Header.Set("Cookie", COOKIE_NAME+"="+validToken)

	respDelete, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respDelete.StatusCode)

	cookies := respDelete.Cookies()
	var cookie *http.Cookie
	for _, c := range cookies {
		if c.Name == COOKIE_NAME {
			cookie = c
			break
		}
	}

	assert.NotNil(t, cookie, "Expected cookie to be found in response")
	assert.Equal(t, "", cookie.Value, "Expected cookie to be deleted (empty value)")
}