package handler_test

import (
	"api-upload-photos/src/infrastructure/controller/handler"
	"api-upload-photos/src/infrastructure/dto"
	"testing"
)

type testCase struct {
	desc    string
	dto     *dto.DTOUser
	expects string
}

func TestProcessUser(t *testing.T) {
	cases := loadTestCases()

	for _, testCase := range cases {
		t.Run(testCase.desc, func(t *testing.T) {
			err := handler.ProcessUser(testCase.dto)
			if err != nil {
				if err.Message != testCase.expects {
					t.Errorf("%s: expected error %s, got %s", testCase.desc, testCase.expects, err.Message)
				}
			} else if testCase.expects != "" {
				t.Errorf("%s: expected error %s, got nil", testCase.desc, testCase.expects)
			}
		})
	}

}

func loadTestCases() []testCase {
	return []testCase{
		{
			desc: "Valid user",
			dto: &dto.DTOUser{
				Username:  "validuser",
				Password:  "ValidPass123!",
				Email:     "validuser@example.com",
				Lastname:  "Doe",
				Firstname: "John",
			},
			expects: "",
		},
		{
			desc: "Password too short",
			dto: &dto.DTOUser{
				Username:  "user1",
				Password:  "Short1!",
				Email:     "user1@example.com",
				Lastname:  "Doe",
				Firstname: "John",
			},
			expects: "La contraseña tiene que tener al menos 8 carácteres",
		},
		{
			desc: "Password missing uppercase",
			dto: &dto.DTOUser{
				Username:  "user2",
				Password:  "lowercase123!",
				Email:     "user2@example.com",
				Lastname:  "Doe",
				Firstname: "John",
			},
			expects: "La contraseña tiene que tener al menos una mayúscula",
		},
		{
			desc: "Password missing special character",
			dto: &dto.DTOUser{
				Username:  "user3",
				Password:  "NoSpecial123",
				Email:     "user3@example.com",
				Lastname:  "Doe",
				Firstname: "John",
			},
			expects: "La contraseña tiene que tener al menos un carácter especial",
		},
		{
			desc: "Invalid email",
			dto: &dto.DTOUser{
				Username:  "user4",
				Password:  "ValidPass123!",
				Email:     "invalid-email",
				Lastname:  "Doe",
				Firstname: "John",
			},
			expects: "El email no es correcto",
		},
	}
}
