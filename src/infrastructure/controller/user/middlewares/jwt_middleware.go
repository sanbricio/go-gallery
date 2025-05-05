package userMiddleware

import (
	"go-gallery/src/commons/exception"
	"go-gallery/src/infrastructure/auth"
	userDTO "go-gallery/src/infrastructure/dto/user"
	log "go-gallery/src/infrastructure/logger"
	userService "go-gallery/src/service/user"
	"time"

	"github.com/gofiber/fiber/v2"
)

var logger log.Logger

const (
	COOKIE_NAME          string        = "auth_token"
	JWT_EXPIRATION_HOURS time.Duration = 2 * time.Hour
)

type JWTMiddleware struct {
	tokenManager auth.TokenManager
}

func NewJWTMiddleware(tokenManager auth.TokenManager) *JWTMiddleware {
	logger = log.Instance()
	return &JWTMiddleware{tokenManager: tokenManager}
}

// Middleware to validate the JWT cookie
func (auth *JWTMiddleware) Handler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Get the JWT token from the cookie
		cookie := ctx.Cookies(COOKIE_NAME)
		if cookie == "" {
			logger.Error("No active session found")
			return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "No active session found"))
		}

		// Get claims from the token
		claims, err := auth.tokenManager.ValidateToken(cookie)
		if err != nil {
			logger.Error("Failed to parse JWT claims: " + err.Message)
			return ctx.Status(err.Status).JSON(err)
		}

		// Check expiration and renew the token if there are less than 10 minutes remaining
		if claims.Expiration-time.Now().Unix() < 600 {
			newToken, err := auth.tokenManager.CreateToken(claims.Username, claims.Email)
			if err != nil {
				logger.Error("Failed to create a new JWT token: " + err.Message)
				return ctx.Status(fiber.StatusInternalServerError).JSON(err)
			}
			auth.createCookie(ctx, newToken)
		}

		// Save the user in the context
		ctx.Locals("user", claims)

		return ctx.Next()
	}
}

func (auth *JWTMiddleware) CreateJWTToken(ctx *fiber.Ctx, username, email string) *exception.ApiException {
	t, err := auth.tokenManager.CreateToken(username, email)
	if err != nil {
		logger.Error("Error creating JWT token: " + err.Message)
		return err
	}
	auth.createCookie(ctx, t)

	return nil
}

func (auth *JWTMiddleware) DeleteAuthCookie(ctx *fiber.Ctx) {
	// Delete the cookie
	ctx.Cookie(&fiber.Cookie{
		Name:     COOKIE_NAME,
		Value:    "",
		MaxAge:   0, // Expires immediately
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	logger.Info("Auth cookie deleted successfully")
}

func ValidateUserClaims(ctx *fiber.Ctx, userService *userService.UserService) (*userDTO.JwtClaimsDTO, *exception.ApiException) {
	claims, ok := ctx.Locals("user").(*userDTO.JwtClaimsDTO)
	if !ok {
		logger.Error("User not authenticated")
		return nil, exception.NewApiException(fiber.StatusUnauthorized, "User not authenticated")
	}

	_, errUser := userService.FindAndCheckJWT(claims)
	if errUser != nil {
		logger.Error("User validation failed: " + errUser.Message)
		return nil, errUser
	}

	logger.Info("User claims validated successfully")
	return claims, nil
}

func (auth *JWTMiddleware) createCookie(ctx *fiber.Ctx, token string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     COOKIE_NAME,
		Value:    token,
		Expires:  time.Now().Add(JWT_EXPIRATION_HOURS),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	logger.Info("Auth cookie created successfully")
}
