package user_repository

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

const UserMongoDBRepositoryKey = "UserMongoDBRepository"

const (
	USER_COLLECTION = "User"
	USERNAME        = "username"
	EMAIL           = "email"
)

type UserMongoDBRepository struct {
	mongo *mongo.Collection
}

func NewUserMongoDBRepository(args map[string]string) UserRepository {
	urlConnection := args["MONGODB_URL_CONNECTION"]
	databaseName := args["MONGODB_DATABASE"]

	db := connect(urlConnection, databaseName)

	repo := &UserMongoDBRepository{
		mongo: db.Collection(USER_COLLECTION),
	}
	return repo
}

func connect(urlConnection, databaseName string) *mongo.Database {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(urlConnection))
	if err != nil {
		panic(fmt.Sprintf("No se ha podido conectar a MongoDB: %s", err.Error()))
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		panic(fmt.Sprintf("No se ha podido hacer ping a MongoDB: %s", err.Error()))
	}

	return client.Database(databaseName)
}

func (r *UserMongoDBRepository) Find(dtoUserFind *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
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

func (r *UserMongoDBRepository) FindAndCheckJWT(claims *dto.DTOClaimsJwt) (*dto.DTOUser, *exception.ApiException) {
	filter := bson.M{USERNAME: claims.Username}
	user, err := r.find(filter)
	if err != nil {
		return nil, err
	}

	if user[0].GetEmail() != claims.Email ||
		user[0].GetUsername() != claims.Username {
		return nil, exception.NewApiException(403, "Los datos proporcionados no coinciden con el usuario autenticado")
	}

	dto := dto.FromUser(user[0])

	return dto, nil
}

func (r *UserMongoDBRepository) Insert(dtoInsertUser *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
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

func (r *UserMongoDBRepository) Update(dtoUpdateUser *dto.DTOUser) (int64, *exception.ApiException) {
	filter := bson.M{"username": dtoUpdateUser.Username}
	_, err := r.find(filter)
	if err != nil {
		return 0, err
	}

	// Construir solo los campos que deben actualizarse
	updateFields := bson.M{}

	if dtoUpdateUser.Email != "" {
		updateFields["email"] = dtoUpdateUser.Email
	}
	if dtoUpdateUser.Firstname != "" {
		updateFields["firstname"] = dtoUpdateUser.Firstname
	}
	if dtoUpdateUser.Lastname != "" {
		updateFields["lastname"] = dtoUpdateUser.Lastname
	}
	if dtoUpdateUser.Password != "" {
		updateFields["password"] = dtoUpdateUser.Password
	}

	if len(updateFields) == 0 {
		return 0, exception.NewApiException(400, "No hay datos para actualizar")
	}

	update := bson.M{"$set": updateFields}
	// Actualizamos el usuario
	result, errUpdate := r.mongo.UpdateOne(context.Background(), filter, update)
	if errUpdate != nil {
		return 0, exception.NewApiException(500, "Error al actualizar el usuario en la base de datos")
	}

	if result.ModifiedCount == 0 {
		return 0, exception.NewApiException(404, "No se ha actualizado ningún usuario")
	}

	return result.ModifiedCount, nil
}

func (r *UserMongoDBRepository) Delete(dtoDeleteUser *dto.DTOUser) (int64, *exception.ApiException) {
	filter := bson.M{USERNAME: dtoDeleteUser.Username}
	user, err := r.find(filter)
	if err != nil {
		return 0, err
	}

	errPassword := user[0].CheckPasswordIntegrity(dtoDeleteUser.Password)
	if errPassword != nil {
		return 0, exception.NewApiException(404, "Contraseña incorrecta")
	}

	result, errDelete := r.mongo.DeleteOne(context.Background(), filter)
	if errDelete != nil {
		return 0, exception.NewApiException(500, "Error al eliminar el usuario")
	}

	deleteCount := result.DeletedCount

	if deleteCount == 0 {
		return 0, exception.NewApiException(404, "No se ha eliminado ningún usuario")
	}

	return deleteCount, nil
}

func (r *UserMongoDBRepository) checkUserIsCreated(dtoInsertUser *dto.DTOUser) *exception.ApiException {
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

func (r *UserMongoDBRepository) find(filter bson.M) ([]*entity.User, *exception.ApiException) {
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
