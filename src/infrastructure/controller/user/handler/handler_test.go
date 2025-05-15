package userHandler

import (
	"testing"

	userDTO "go-gallery/src/infrastructure/dto/user"

	"github.com/stretchr/testify/assert"
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

			if testCase.expects == "" {
				assert.Nil(t, err, testCase.desc)
			} else {
				assert.NotNil(t, err, testCase.desc)
				assert.Equal(t, testCase.expects, err.Message, testCase.desc)
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
			expects: "The password must be at least 8 characters long",
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
			expects: "The password must contain at least one uppercase letter",
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
			expects: "The password must contain at least one special character",
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
			expects: "The email is not valid",
		},
	}
}
