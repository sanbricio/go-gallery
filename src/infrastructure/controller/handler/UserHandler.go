package handler

import (
	"go-gallery/src/commons/exception"
	"net/http"
	"regexp"
)

func ProcessUser(password, email string) *exception.ApiException {
	if password != "" {
		err := validatePassword(password)
		if err != nil {
			return err
		}
	}

	if email != "" {
		err := validateEmail(email)
		if err != nil {
			return err
		}
	}

	return nil
}

func validatePassword(password string) *exception.ApiException {
	if len(password) < 8 {
		return exception.NewApiException(http.StatusBadRequest, "La contraseña tiene que tener al menos 8 carácteres")
	}

	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUppercase {
		return exception.NewApiException(http.StatusBadRequest, "La contraseña tiene que tener al menos una mayúscula")
	}

	hasSpecialChar := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	if !hasSpecialChar {
		return exception.NewApiException(http.StatusBadRequest, "La contraseña tiene que tener al menos un carácter especial")
	}

	return nil
}

func validateEmail(email string) *exception.ApiException {
	const emailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailPattern)
	if !re.MatchString(email) {
		return exception.NewApiException(http.StatusBadRequest, "El email no es correcto")
	}
	return nil
}
