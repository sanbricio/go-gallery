package repository

import (
	"api-upload-photos/src/commons/exception"
	entity "api-upload-photos/src/domain/entities"
	"api-upload-photos/src/domain/entities/builder"
	"api-upload-photos/src/infrastructure/dto"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	USER_COLLECTION = "User"
	USERNAME        = "username"
	EMAIL           = "email"
)

type RepositoryUserMongoDB struct {
	mongo *mongo.Collection
}

func NewRepositoryMongoDB(urlConnection, databaseName string) (*RepositoryUserMongoDB, *exception.ConnectionException) {
	db, err := connect(urlConnection, databaseName)
	if err != nil {
		return nil, err
	}

	repo := &RepositoryUserMongoDB{
		mongo: db.Collection(USER_COLLECTION),
	}
	return repo, nil
}

func connect(urlConnection, databaseName string) (*mongo.Database, *exception.ConnectionException) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(urlConnection))
	if err != nil {
		return nil, exception.NewConnectionException(fmt.Sprintf("No se ha podido conectar a MongoDB: %s", err.Error()), err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, exception.NewConnectionException(fmt.Sprintf("No se ha podido hacer ping a MongoDB: %s", err.Error()), err)
	}

	database := client.Database(databaseName)
	return database, nil
}

func (r *RepositoryUserMongoDB) Find(dtoUserFind *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	filter := bson.M{USERNAME: dtoUserFind.Username}
	user, err := r.find(filter)
	if err != nil {
		return nil, err
	}

	errPasword := user[0].CheckPasswordIntegrity(dtoUserFind.Password)
	if errPasword != nil {
		return nil, exception.NewApiException(404, "Contraseña incorrecta")
	}

	dto := dto.FromUser(user[0])

	return dto, nil
}

func (r *RepositoryUserMongoDB) FindAndCheckJWT(dtoUserFind *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	filter := bson.M{USERNAME: dtoUserFind.Username}
	user, err := r.find(filter)
	if err != nil {
		return nil, err
	}

	if user[0].GetEmail() != dtoUserFind.Email ||
		user[0].GetFirstname() != dtoUserFind.Firstname ||
		user[0].GetUsername() != dtoUserFind.Username {
		return nil, exception.NewApiException(403, "Los datos proporcionados no coinciden con el usuario autenticado")
	}

	dto := dto.FromUser(user[0])

	return dto, nil
}

func (r *RepositoryUserMongoDB) Insert(dtoInsertUser *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	err := r.checkUserIsCreated(dtoInsertUser)
	if err != nil {
		return nil, err
	}

	user, errBuilder := builder.NewUserBuilder().
		FromDTO(dtoInsertUser).
		Build()
	if errBuilder != nil {
		return nil, exception.NewApiException(500, errBuilder.Error())
	}

	dto := dto.FromUser(user)

	_, errInsert := r.mongo.InsertOne(context.Background(), dto)
	if errInsert != nil {
		return nil, exception.NewApiException(500, "Error al insertar el documento")
	}

	return dto, nil
}

func (r *RepositoryUserMongoDB) Update(dtoUpdateUser *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	filter := bson.M{USERNAME: dtoUpdateUser.Username}
	_, err := r.find(filter)
	if err != nil {
		return nil, err
	}

	user, errBuilder := builder.NewUserBuilder().
		FromDTO(dtoUpdateUser).
		Build()
	if errBuilder != nil {
		return nil, exception.NewApiException(500, errBuilder.Error())
	}

	update := bson.M{
		"$set": bson.M{
			"email":     user.GetEmail(),
			"firstname": user.GetFirstname(),
			"lastname":  user.GetLastname(),
			"password":  user.GetPassword(),
		},
	}

	_, errUpdate := r.mongo.UpdateOne(context.Background(), filter, update)
	if errUpdate != nil {
		return nil, exception.NewApiException(500, "Error al actualizar el usuario en la base de datos")
	}

	updatedDTO := dto.FromUser(user)
	return updatedDTO, nil
}

func (r *RepositoryUserMongoDB) Delete(dtoDeleteUser *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	filter := bson.M{USERNAME: dtoDeleteUser.Username}
	user, err := r.find(filter)
	if err != nil {
		return nil, err
	}

	errPassword := user[0].CheckPasswordIntegrity(dtoDeleteUser.Password)
	if errPassword != nil {
		return nil, exception.NewApiException(404, "Contraseña incorrecta")
	}

	_, errDelete := r.mongo.DeleteOne(context.Background(), filter)
	if errDelete != nil {
		return nil, exception.NewApiException(500, "Error al eliminar el usuario")
	}

	deletedUserDTO := dto.FromUser(user[0])
	return deletedUserDTO, nil
}

func (r *RepositoryUserMongoDB) checkUserIsCreated(dtoInsertUser *dto.DTOUser) *exception.ApiException {
	filter := bson.M{
		"$or": []bson.M{
			{EMAIL: dtoInsertUser.Email},
			{USERNAME: dtoInsertUser.Username},
		},
	}

	users, err := r.find(filter)
	if err != nil && err.Status != 404 {
		return err
	}

	for _, user := range users {
		if user.GetEmail() == dtoInsertUser.Email {
			return exception.NewApiException(400, "El correo electrónico ya está registrado")
		}
		if user.GetUsername() == dtoInsertUser.Username {
			return exception.NewApiException(400, "El usuario ya está registrado")
		}
	}

	return nil
}

func (r *RepositoryUserMongoDB) find(filter bson.M) ([]*entity.User, *exception.ApiException) {
	cursor, err := r.mongo.Find(context.Background(), filter)
	if err != nil {
		return nil, exception.NewApiException(500, "Error al buscar usuarios")
	}
	defer cursor.Close(context.Background())

	var users []*dto.DTOUser
	if err := cursor.All(context.Background(), &users); err != nil {
		return nil, exception.NewApiException(500, "Error al decodificar usuarios")
	}

	if len(users) == 0 {
		return nil, exception.NewApiException(404, "Usuario no encontrado")
	}

	var userEntities []*entity.User
	for _, userDTO := range users {
		user, errBuilder := builder.NewUserBuilder().
			FromDTO(userDTO).
			Build()
		if errBuilder != nil {
			return nil, exception.NewApiException(500, errBuilder.Error())
		}
		userEntities = append(userEntities, user)
	}

	return userEntities, nil
}
