package userHandler

import (
	userDTO "go-gallery/src/infrastructure/dto/user"
	"testing"
)

type testCaseUserHandler struct {
	desc    string
	dto     *userDTO.UserDTO
	expects string
}

func TestProcessUser(t *testing.T) {
	cases := loadTestCasesUserHandler()

	for _, testCase := range cases {
		t.Run(testCase.desc, func(t *testing.T) {
			err := ProcessUser(testCase.dto.Password, testCase.dto.Email)
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

func loadTestCasesUserHandler() []testCaseUserHandler {
	return []testCaseUserHandler{
		{
			desc: "Valid user",
			dto: &userDTO.UserDTO{
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
			dto: &userDTO.UserDTO{
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
			dto: &userDTO.UserDTO{
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
			dto: &userDTO.UserDTO{
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
			dto: &userDTO.UserDTO{
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
