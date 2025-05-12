package userHandler

import (
	"go-gallery/src/commons/exception"
	"net/http"
	"regexp"
)

func ProcessUser(password, email string) *exception.ApiException {
	if password != "" {
		err := ValidatePassword(password)
		if err != nil {
			return err
		}
	}

	if email != "" {
		err := ValidateEmail(email)
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidatePassword(password string) *exception.ApiException {
	if len(password) < 8 {
		return exception.NewApiException(http.StatusBadRequest, "The password must be at least 8 characters long")
	}

	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUppercase {
		return exception.NewApiException(http.StatusBadRequest, "The password must contain at least one uppercase letter")
	}

	hasSpecialChar := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	if !hasSpecialChar {
		return exception.NewApiException(http.StatusBadRequest, "The password must contain at least one special character")
	}

	return nil
}

func ValidateEmail(email string) *exception.ApiException {
	const emailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailPattern)
	if !re.MatchString(email) {
		return exception.NewApiException(http.StatusBadRequest, "The email is not valid")
	}
	return nil
}
