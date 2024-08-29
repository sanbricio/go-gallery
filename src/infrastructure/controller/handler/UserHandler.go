package handler

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/dto"
	"regexp"
)

func ProcessUser(dto *dto.DTORegisterRequest) *exception.ApiException {
	if dto.Username == "" {
		return exception.NewApiException(404, "El usuario no puede estar vacio")
	}

	errPassword := validatePassword(dto.Password)
	if errPassword != nil {
		return errPassword
	}

	errEmail := validateEmail(dto.Email)
	if errEmail != nil {
		return errPassword
	}

	return nil
}

func validatePassword(password string) *exception.ApiException {
	if len(password) < 8 {
		return exception.NewApiException(404, "La contraseña tiene que tener al menos 8 carácteres")
	}

	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUppercase {
		return exception.NewApiException(404, "La contraseña tiene que tener al menos una mayúscula")
	}

	hasSpecialChar := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	if !hasSpecialChar {
		return exception.NewApiException(404, "La contraseña tiene que tener al menos un carácter especial")
	}

	return nil
}

func validateEmail(email string) *exception.ApiException {
	const emailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailPattern)
	if !re.MatchString(email) {
		return exception.NewApiException(404, "El email no es correcto")
	}
	return nil
}
