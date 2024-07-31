package dto

import entity "api-upload-photos/src/domain/entities"

type DTOUser struct {
	Username  string `json:"username" bson:"username"`
	Password  string `json:"password" bson:"pasword"`
	Email     string `json:"email" bson:"email"`
	Lastname  string `json:"lastname" bson:"lastname"`
	Firstname string `json:"firstname" bson:"firstname"`
}

func FromUser(user *entity.User) *DTOUser {
	return &DTOUser{
		Username:  user.GetUsername(),
		Password:  user.GetPassword(),
		Email:     user.GetEmail(),
		Lastname:  user.GetLastname(),
		Firstname: user.GetFirstname(),
	}
}

func (dto *DTOUser) AsUserEntity() *entity.User {
	return entity.NewUserFromDTO(dto.Username, dto.Password, dto.Email, dto.Firstname, dto.Lastname)
}
