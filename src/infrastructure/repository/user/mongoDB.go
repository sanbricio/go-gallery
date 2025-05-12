package userRepository

import (
	"context"
	"fmt"
	"go-gallery/src/commons/exception"
	userBuilder "go-gallery/src/domain/entities/builder/user"
	userEntity "go-gallery/src/domain/entities/user"
	userDTO "go-gallery/src/infrastructure/dto/user"
	log "go-gallery/src/infrastructure/logger"

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

	logger = log.Instance()

	db := connect(urlConnection, databaseName)

	repo := &UserMongoDBRepository{
		mongo: db.Collection(USER_COLLECTION),
	}
	return repo
}

func connect(urlConnection, databaseName string) *mongo.Database {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(urlConnection))
	if err != nil {
		panicMessage := fmt.Sprintf("Unable to connect to MongoDB: %s", err.Error())
		logger.Panic(panicMessage)
		panic(panicMessage)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		panicMessage := fmt.Sprintf("Unable to ping MongoDB: %s", err.Error())
		logger.Panic(panicMessage)
		panic(panicMessage)
	}

	return client.Database(databaseName)
}

func (r *UserMongoDBRepository) Find(dtoUserFind *userDTO.LoginRequestDTO) (*userDTO.UserDTO, *exception.ApiException) {
	logger.Info(fmt.Sprintf("Searching for user: %s", dtoUserFind.Username))
	filter := bson.M{USERNAME: dtoUserFind.Username}
	user, err := r.find(filter)
	if err != nil {
		logger.Warning(fmt.Sprintf("User not found: %s", dtoUserFind.Username))
		return nil, err
	}

	errPassword := user[0].CheckPasswordIntegrity(dtoUserFind.Password)
	if errPassword != nil {
		logger.Warning(fmt.Sprintf("Incorrect password for user: %s", dtoUserFind.Username))
		return nil, exception.NewApiException(404, "Incorrect password")
	}

	dto := userDTO.FromUser(user[0])

	logger.Info(fmt.Sprintf("User found: %s", user[0].GetUsername()))

	return dto, nil
}

func (r *UserMongoDBRepository) FindByEmail(email string) (*userDTO.UserDTO, *exception.ApiException) {
	panic("method not implemented FindByEmail in UserMongoDBRepository")
}

func (r *UserMongoDBRepository) FindAndCheckJWT(claims *userDTO.JwtClaimsDTO) (*userDTO.UserDTO, *exception.ApiException) {
	logger.Info(fmt.Sprintf("Verifying JWT for user: %s", claims.Username))

	filter := bson.M{USERNAME: claims.Username}
	user, err := r.find(filter)
	if err != nil {
		logger.Warning(fmt.Sprintf("User not found while verifying JWT: %s", claims.Username))
		return nil, err
	}

	if user[0].GetEmail() != claims.Email ||
		user[0].GetUsername() != claims.Username {
		logger.Warning(fmt.Sprintf(
			"Data mismatch in JWT verification. Username: expected %s / obtained %s, Email: expected %s / obtained %s",
			claims.Username, user[0].GetUsername(), claims.Email, user[0].GetEmail(),
		))
		return nil, exception.NewApiException(403, "The provided data does not match the authenticated user")
	}

	dto := userDTO.FromUser(user[0])

	logger.Info(fmt.Sprintf("JWT successfully verified for user: %s", claims.Username))

	return dto, nil
}

func (r *UserMongoDBRepository) Insert(dtoInsertUser *userDTO.UserDTO) (*userDTO.UserDTO, *exception.ApiException) {
	logger.Info(fmt.Sprintf("Attempting to register user: %s", dtoInsertUser.Username))
	err := r.checkUserIsCreated(dtoInsertUser)
	if err != nil {
		logger.Warning(fmt.Sprintf("User already exists or verification error: %s", err.Message))
		return nil, err
	}

	user, errBuilder := userBuilder.NewUserBuilder().
		FromDTO(dtoInsertUser).
		Build()
	if errBuilder != nil {
		logger.Error(fmt.Sprintf("Error building user: %s", errBuilder.Error()))
		return nil, exception.NewApiException(500, errBuilder.Error())
	}

	dto := userDTO.FromUser(user)

	_, errInsert := r.mongo.InsertOne(context.Background(), dto)
	if errInsert != nil {
		logger.Error(fmt.Sprintf("Error inserting user %s: %s", user.GetUsername(), errInsert.Error()))
		return nil, exception.NewApiException(500, "Error inserting document")
	}

	logger.Info(fmt.Sprintf("User successfully inserted: %s", user.GetUsername()))
	return dto, nil
}

func (r *UserMongoDBRepository) Update(dtoUpdateUser *userDTO.UserDTO) (int64, *exception.ApiException) {
	logger.Info(fmt.Sprintf("Attempting to update user: %s", dtoUpdateUser.Username))

	filter := bson.M{"username": dtoUpdateUser.Username}
	_, err := r.find(filter)
	if err != nil {
		logger.Warning(fmt.Sprintf("User not found: %s", dtoUpdateUser.Username))
		return 0, err
	}

	// Build only the fields that need to be updated
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
		logger.Warning(fmt.Sprintf("No data to update for user: %s", dtoUpdateUser.Username))
		return 0, exception.NewApiException(400, "No data to update")
	}

	update := bson.M{"$set": updateFields}
	// Update the user
	result, errUpdate := r.mongo.UpdateOne(context.Background(), filter, update)
	if errUpdate != nil {
		logger.Error(fmt.Sprintf("Error updating user %s: %s", dtoUpdateUser.Username, errUpdate.Error()))
		return 0, exception.NewApiException(500, "Error updating the user in the database")
	}

	if result.ModifiedCount == 0 {
		logger.Warning(fmt.Sprintf("User not found to update: %s", dtoUpdateUser.Username))
		return 0, exception.NewApiException(404, "No user updated")
	}

	return result.ModifiedCount, nil
}

func (r *UserMongoDBRepository) Delete(dtoDeleteUser *userDTO.UserDTO) (int64, *exception.ApiException) {
	logger.Info(fmt.Sprintf("Attempting to delete user: %s", dtoDeleteUser.Username))

	filter := bson.M{USERNAME: dtoDeleteUser.Username}
	user, err := r.find(filter)
	if err != nil {
		logger.Warning(fmt.Sprintf("Error searching for user to delete %s: %s", dtoDeleteUser.Username, err.Message))
		return 0, err
	}

	errPassword := user[0].CheckPasswordIntegrity(dtoDeleteUser.Password)
	if errPassword != nil {
		return 0, exception.NewApiException(404, "Incorrect password")
	}

	result, errDelete := r.mongo.DeleteOne(context.Background(), filter)
	if errDelete != nil {
		logger.Error(fmt.Sprintf("Error deleting user %s: %s", dtoDeleteUser.Username, errDelete.Error()))
		return 0, exception.NewApiException(500, "Error deleting the user")
	}

	deleteCount := result.DeletedCount

	if deleteCount == 0 {
		logger.Warning(fmt.Sprintf("User not found to delete: %s", dtoDeleteUser.Username))
		return 0, exception.NewApiException(404, "No user deleted")
	}

	return deleteCount, nil
}

func (r *UserMongoDBRepository) checkUserIsCreated(dtoInsertUser *userDTO.UserDTO) *exception.ApiException {
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
			return exception.NewApiException(400, "Email is already registered")
		}
		if user.GetUsername() == dtoInsertUser.Username {
			return exception.NewApiException(400, "Username is already registered")
		}
	}

	return nil
}

func (r *UserMongoDBRepository) find(filter bson.M) ([]*userEntity.User, *exception.ApiException) {
	cursor, err := r.mongo.Find(context.Background(), filter)
	if err != nil {
		return nil, exception.NewApiException(500, "Error searching for users")
	}
	defer cursor.Close(context.Background())

	var users []*userDTO.UserDTO
	if err := cursor.All(context.Background(), &users); err != nil {
		return nil, exception.NewApiException(500, "Error decoding users")
	}

	if len(users) == 0 {
		return nil, exception.NewApiException(404, "User not found")
	}

	var userEntities []*userEntity.User
	for _, userDTO := range users {
		user, errBuilder := userBuilder.NewUserBuilder().
			FromDTO(userDTO).
			Build()
		if errBuilder != nil {
			return nil, exception.NewApiException(500, errBuilder.Error())
		}
		userEntities = append(userEntities, user)
	}

	return userEntities, nil
}
