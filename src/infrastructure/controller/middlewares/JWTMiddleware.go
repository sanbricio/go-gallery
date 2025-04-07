package middlewares

import (
	"go-gallery/src/commons/exception"
	"go-gallery/src/infrastructure/dto"
	"time"

	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v5"
)

const (
	COOKIE_NAME          string        = "auth_token"
	JWT_EXPIRATION_HOURS time.Duration = 2 * time.Hour
)

type JWTMiddleware struct {
	secret string
}

func NewJWTMiddleware(secret string) *JWTMiddleware {
	return &JWTMiddleware{secret: secret}
}

// Middleware para validar la cookie JWT
func (auth *JWTMiddleware) Handler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Obtener la cookie con el token JWT
		cookie := ctx.Cookies(COOKIE_NAME)
		if cookie == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "No se ha encontrado una sesión activa"))
		}

		// Obtener claims del token
		claims, err := auth.getJWTClaimsFromCookie(cookie)
		if err != nil {
			return ctx.Status(err.Status).JSON(err)
		}

		// Verificar expiración y renovar si quedan menos de 10 minutos
		if claims.Expiration-time.Now().Unix() < 600 {
			newToken, err := auth.createJwtToken(claims.Username, claims.Email)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(err)
			}
			auth.createCookie(ctx, newToken)
		}

		// Guardar el usuario en el contexto
		ctx.Locals("user", claims)

		return ctx.Next()
	}
}

func (auth *JWTMiddleware) getJWTClaimsFromCookie(cookie string) (*dto.DTOClaimsJwt, *exception.ApiException) {
	token, err := jtoken.Parse(cookie, func(token *jtoken.Token) (any, error) {
		return []byte(auth.secret), nil
	})

	if err != nil || !token.Valid {
		return nil, exception.NewApiException(fiber.StatusUnauthorized, "Token inválido")
	}

	claims, ok := token.Claims.(jtoken.MapClaims)
	if !ok {
		return nil, exception.NewApiException(fiber.StatusInternalServerError, "Error al obtener las claims")
	}

	username, ok := claims["username"].(string)
	email, okEmail := claims["email"].(string)
	iat, okFirstname := claims["iat"].(float64)
	exp, okExp := claims["exp"].(float64)

	if !ok || !okEmail || !okFirstname || !okExp {
		return nil, exception.NewApiException(fiber.StatusInternalServerError, "Error en claims JWT")
	}

	return &dto.DTOClaimsJwt{
		Username:   username,
		Email:      email,
		IssuedAt:   int64(iat),
		Expiration: int64(exp),
	}, nil
}

func (auth *JWTMiddleware) CreateJWTToken(ctx *fiber.Ctx, username, email string) *exception.ApiException {
	t, err := auth.createJwtToken(username, email)
	if err != nil {
		return err
	}
	auth.createCookie(ctx, t)

	return nil
}

func (auth *JWTMiddleware) createJwtToken(username, email string) (string, *exception.ApiException) {
	// Crear las claims del JWT, incluyendo el usuario y la expiración
	claims := jtoken.MapClaims{
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(JWT_EXPIRATION_HOURS).Unix(), // Fecha de expiración del token
		"iat":      time.Now().Unix(),                           // Fecha de emisión del token
	}

	// Creamos el token JWT
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)

	// Firmar el token
	t, err := token.SignedString([]byte(auth.secret))
	if err != nil {
		return "", exception.NewApiException(500, "Error al crear el token JWT")
	}

	return t, nil
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
}

func (auth *JWTMiddleware) DeleteAuthCookie(ctx *fiber.Ctx) {
	// Eliminar la cookie
	ctx.Cookie(&fiber.Cookie{
		Name:     COOKIE_NAME,
		Value:    "",
		MaxAge:   0, // Expira inmediatamente
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})
}
