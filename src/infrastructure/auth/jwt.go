package auth

import (
	"go-gallery/src/commons/exception"
	userDTO "go-gallery/src/infrastructure/dto/user"
	log "go-gallery/src/infrastructure/logger"
	"time"

	jtoken "github.com/golang-jwt/jwt/v5"
)

const (
	JWT_EXPIRATION_HOURS time.Duration = 2 * time.Hour
)

type JWTTokenManager struct {
	secret string
}

var logger log.Logger

func NewJWTTokenManager(secret string) *JWTTokenManager {
	logger = log.Instance()
	return &JWTTokenManager{secret: secret}
}

func (j *JWTTokenManager) CreateToken(username, email string) (string, *exception.ApiException) {
	// Create the JWT claims, including the user and expiration
	claims := jtoken.MapClaims{
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(JWT_EXPIRATION_HOURS).Unix(), // Expiration date of the token
		"iat":      time.Now().Unix(),                           // Issued date of the token
	}

	// Create the JWT token
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)

	// Sign the token
	t, err := token.SignedString([]byte(j.secret))
	if err != nil {
		logger.Error("Failed to sign JWT token: " + err.Error())
		return "", exception.NewApiException(500, "Error creating JWT token")
	}

	logger.Info("JWT token created successfully")
	return t, nil
}

func (j *JWTTokenManager) ValidateToken(tokenString string) (*userDTO.JwtClaimsDTO, *exception.ApiException) {
	token, err := jtoken.Parse(tokenString, func(token *jtoken.Token) (any, error) {
		return []byte(j.secret), nil
	})

	if err != nil || !token.Valid {
		logger.Error("Invalid token")
		return nil, exception.NewApiException(401, "Invalid token")
	}

	claims, ok := token.Claims.(jtoken.MapClaims)
	if !ok {
		logger.Error("Error parsing JWT claims")
		return nil, exception.NewApiException(500, "Error parsing JWT claims")
	}

	username, ok := claims["username"].(string)
	email, okEmail := claims["email"].(string)
	iat, okIat := claims["iat"].(float64)
	exp, okExp := claims["exp"].(float64)

	if !ok || !okEmail || !okIat || !okExp {
		logger.Error("Error in JWT claims")
		return nil, exception.NewApiException(500, "Error in JWT claims")
	}

	return &userDTO.JwtClaimsDTO{
		Username:   username,
		Email:      email,
		IssuedAt:   int64(iat),
		Expiration: int64(exp),
	}, nil
}
