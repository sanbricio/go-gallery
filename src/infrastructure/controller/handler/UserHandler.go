package handler

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/dto"
	"net/http"
	"regexp"
)

func ProcessUser(dto *dto.DTOUser) *exception.ApiException {
	if dto.Password != "" {
		err := validatePassword(dto.Password)
		if err != nil {
			return err
		}
	}

	if dto.Email != "" {
		err := validateEmail(dto.Email)
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
