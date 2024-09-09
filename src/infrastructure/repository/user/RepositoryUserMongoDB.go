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
	client *mongo.Database
}

func NewRepositoryMongoDB(urlConnection, databaseName string) (*RepositoryUserMongoDB, *exception.ConnectionException) {
	db, err := connect(urlConnection, databaseName)
	if err != nil {
		return nil, err
	}

	repo := &RepositoryUserMongoDB{
		client: db,
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

func (r *RepositoryUserMongoDB) Find(dtoLoginRequest *dto.DTOLoginRequest) (*dto.DTOUser, *exception.ApiException) {
	filter := bson.M{USERNAME: dtoLoginRequest.Username}
	user, err := r.find(filter)
	if err != nil {
		return nil, err
	}

	errPasword := user.CheckPasswordIntegrity(dtoLoginRequest.Password)
	if errPasword != nil {
		return nil, exception.NewApiException(404, "Contraseña incorrecta")
	}

	dto := dto.FromUser(user)

	return dto, nil
}

func (r *RepositoryUserMongoDB) find(filter bson.M) (*entity.User, *exception.ApiException) {
	collection := r.client.Collection(USER_COLLECTION)

	var dto *dto.DTOUser
	err := collection.FindOne(context.Background(), filter).Decode(&dto)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, exception.NewApiException(404, "No se ha encontrado el usuario")
		}
		return nil, exception.NewApiException(500, fmt.Sprintf("Error al buscar el usuario: %s", err.Error()))
	}

	user, err := builder.NewUserBuilder().
		FromDTO(dto).
		Build()

	return user, nil
}

func (r *RepositoryUserMongoDB) Insert(dtoRegister *dto.DTORegisterRequest) (*dto.DTOUser, *exception.ApiException) {
	collection := r.client.Collection(USER_COLLECTION)

	err := r.checkUserIsCreated(dtoRegister)
	if err != nil {
		return nil, err
	}

	user, errBuilder := builder.NewUserBuilder().
		FromDTORegister(dtoRegister).
		Build()
	if errBuilder != nil {
		return nil, nil
	}

	dto := dto.FromUser(user)

	_, errInsert := collection.InsertOne(context.Background(), dto)
	if errInsert != nil {
		return nil, exception.NewApiException(500, "Error al insertar el documento")
	}

	return dto, nil
}

func (r *RepositoryUserMongoDB) checkUserIsCreated(dtoRegister *dto.DTORegisterRequest) *exception.ApiException {
	//TODO Creo que se puede quitar un find, reestructurar el r.find para que devuevla un findAll
	// Cuando tenga el array de documents validar si el email es igual al del registro o el usuario
	existingEmail, err := r.find(bson.M{EMAIL: dtoRegister.Email})
	if err != nil && err.Status != 404 {
		return err
	}

	if existingEmail != nil {
		return exception.NewApiException(400, "El correo electrónico ya está registrado")
	}

	existingUser, err := r.find(bson.M{USERNAME: dtoRegister.Username})
	if err != nil && err.Status != 404 {
		return err
	}

	if existingUser != nil {
		return exception.NewApiException(400, "El usuario ya esta registrado")
	}

	return nil
}

// Delete implements IRepositoryUser.
func (r *RepositoryUserMongoDB) Delete() (*dto.DTOUser, *exception.ApiException) {
	panic("unimplemented")
}

// Update implements IRepositoryUser.
func (r *RepositoryUserMongoDB) Update() (*dto.DTOUser, *exception.ApiException) {
	panic("unimplemented")
}
