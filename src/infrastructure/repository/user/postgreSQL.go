package userRepository

import (
	"database/sql"
	"fmt"
	"go-gallery/src/commons/exception"

	userBuilder "go-gallery/src/domain/entities/builder/user"
	userEntity "go-gallery/src/domain/entities/user"

	userDTO "go-gallery/src/infrastructure/dto/user"
	log "go-gallery/src/infrastructure/logger"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const UserPostgreSQLRepositoryKey = "UserPostgreSQLRepository"

var logger log.Logger

type UserPostgreSQLRepository struct {
	db *sql.DB
}

const (
	retries uint = 5
)

func NewUserPostgreSQLRepository(args map[string]string) UserRepository {
	user := args["POSTGRESQL_USER"]
	password := args["POSTGRESQL_PASSWORD"]
	dbName := args["POSTGRESQL_DB"]
	host := args["POSTGRESQL_HOST"]
	port := args["POSTGRESQL_PORT"]

	// Construimos la URL de conexión
	urlConnection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

	logger = log.Instance()

	db, err := sql.Open("postgres", urlConnection)
	if err != nil {
		panic(fmt.Sprintf("Could not connect to PostgreSQL: %s", err.Error()))
	}

	// Comprobamos si la base de datos realmente está disponible
	for i := range retries {
		err = db.Ping()
		if err == nil {
			logger.Info("Successfully connected to PostgreSQL user database")
			break
		}

		// Si hemos llegado al último intento, mostramos un mensaje de error
		if i == retries-1 {
			logger.Error("Could not connect to PostgreSQL after several attempts.")
			panicMessage := fmt.Sprintf("Error trying to connect to PostgreSQL: %s", err.Error())
			logger.Panic(panicMessage)
			panic(panicMessage)
		}

		logger.Warning(fmt.Sprintf("Attempt %d of %d. Retrying...\n", i+1, retries))
		time.Sleep(5 * time.Second)
	}

	// Ejecutar el DDL para crear la tabla si no existe
	ddl, err := os.ReadFile("sql/userTableDDL.sql")
	if err != nil {
		panicMessage := fmt.Sprintf("Could not read the DDL file: %s", err.Error())
		logger.Panic(panicMessage)
		panic(panicMessage)
	}

	_, err = db.Exec(string(ddl))
	if err != nil {
		panicMessage := fmt.Sprintf("Error executing DDL creation: %s", err.Error())
		logger.Panic(panicMessage)
		panic(panicMessage)
	}

	return &UserPostgreSQLRepository{db: db}
}

func (u *UserPostgreSQLRepository) Find(dtoLoginRequest *userDTO.LoginRequestDTO) (*userDTO.UserDTO, *exception.ApiException) {
	logger.Info(fmt.Sprintf("Searching for user: %s", dtoLoginRequest.Username))

	user, err := u.find(dtoLoginRequest.Username, dtoLoginRequest.Password)
	if err != nil {
		logger.Error(fmt.Sprintf("Error searching for user %s: %s", dtoLoginRequest.Username, err.Message))
		return nil, err
	}

	logger.Info(fmt.Sprintf("User found: %s", user.GetUsername()))

	return userDTO.FromUser(user), nil
}

func (u *UserPostgreSQLRepository) find(username, password string) (*userEntity.User, *exception.ApiException) {
	user, err := u.findBy("username", username)
	if err != nil {
		return nil, err
	}

	if err := user.CheckPasswordIntegrity(password); err != nil {
		return nil, exception.NewApiException(400, "Incorrect password")
	}

	return user, nil
}

func (u *UserPostgreSQLRepository) FindByEmail(email string) (*userDTO.UserDTO, *exception.ApiException) {
	logger.Info(fmt.Sprintf("Searching for user by email: %s", email))

	user, err := u.findBy("email", email)
	if err != nil {
		logger.Error(fmt.Sprintf("Error searching for user with email %s: %s", email, err.Message))
		return nil, err
	}

	return userDTO.FromUser(user), nil
}

func (u *UserPostgreSQLRepository) findBy(field, value string) (*userEntity.User, *exception.ApiException) {
	query := fmt.Sprintf("SELECT username, email, firstname, lastname, password FROM users WHERE %s = $1", field)
	row := u.db.QueryRow(query, value)

	userDTO := new(userDTO.UserDTO)
	if err := row.Scan(&userDTO.Username, &userDTO.Email, &userDTO.Firstname, &userDTO.Lastname, &userDTO.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewApiException(404, "User not found")
		}
		return nil, exception.NewApiException(500, "Error retrieving user")
	}

	user, errBuilder := userBuilder.NewUserBuilder().
		FromDTO(userDTO).
		Build()

	if errBuilder != nil {
		return nil, exception.NewApiException(500, errBuilder.Error())
	}

	return user, nil
}

func (u *UserPostgreSQLRepository) FindAndCheckJWT(claims *userDTO.JwtClaimsDTO) (*userDTO.UserDTO, *exception.ApiException) {
	logger.Info(fmt.Sprintf("Verifying JWT for user: %s", claims.Username))

	query := "SELECT username, email FROM users WHERE username = $1"
	row := u.db.QueryRow(query, claims.Username)

	userDTO := new(userDTO.UserDTO)
	if err := row.Scan(&userDTO.Username, &userDTO.Email); err != nil {
		logger.Warning(fmt.Sprintf("User not found when verifying JWT: %s", claims.Username))
		return nil, exception.NewApiException(404, "User not found")
	}

	if userDTO.Username != claims.Username || userDTO.Email != claims.Email {
		logger.Warning(fmt.Sprintf(
			"Data mismatch during JWT verification. Username: expected %s / got %s, Email: expected %s / got %s",
			claims.Username, userDTO.Username, claims.Email, userDTO.Email,
		))
		return nil, exception.NewApiException(403, "The provided data does not match the authenticated user")
	}

	logger.Info(fmt.Sprintf("JWT successfully verified for user: %s", claims.Username))

	return userDTO, nil
}

func (u *UserPostgreSQLRepository) Insert(dtoRegisterRequest *userDTO.UserDTO) (*userDTO.UserDTO, *exception.ApiException) {
	logger.Info(fmt.Sprintf("Attempting to register user: %s", dtoRegisterRequest.Username))

	if err := u.checkUserIsCreated(dtoRegisterRequest); err != nil {
		logger.Warning(fmt.Sprintf("User already exists or verification error: %s", err.Message))
		return nil, err
	}

	user, err := userBuilder.NewUserBuilder().
		FromDTO(dtoRegisterRequest).
		Build()

	if err != nil {
		logger.Error(fmt.Sprintf("Error building user: %s", err.Error()))
		return nil, exception.NewApiException(500, err.Error())
	}

	query := "INSERT INTO users (username, email, firstname, lastname, password) VALUES ($1, $2, $3, $4, $5)"
	_, errDb := u.db.Exec(query, user.GetUsername(), user.GetEmail(), user.GetFirstname(), user.GetLastname(), user.GetPassword())
	if errDb != nil {
		logger.Error(fmt.Sprintf("Error inserting user %s: %s", user.GetUsername(), errDb.Error()))
		return nil, exception.NewApiException(500, "Error inserting user")
	}

	logger.Info(fmt.Sprintf("User successfully inserted: %s", user.GetUsername()))
	return userDTO.FromUser(user), nil
}

func (u *UserPostgreSQLRepository) Update(dtoUpdateUser *userDTO.UserDTO) (int64, *exception.ApiException) {
	logger.Info(fmt.Sprintf("Attempting to update user: %s", dtoUpdateUser.Username))

	query := "UPDATE users SET "
	args := []any{}
	count := 1
	// Dynamically construct the query and arguments for the update
	if dtoUpdateUser.Email != "" {
		query += "email = $" + fmt.Sprint(count) + ", "
		args = append(args, dtoUpdateUser.Email)
		count++
	}
	if dtoUpdateUser.Firstname != "" {
		query += "firstname = $" + fmt.Sprint(count) + ", "
		args = append(args, dtoUpdateUser.Firstname)
		count++
	}
	if dtoUpdateUser.Lastname != "" {
		query += "lastname = $" + fmt.Sprint(count) + ", "
		args = append(args, dtoUpdateUser.Lastname)
		count++
	}
	if dtoUpdateUser.Password != "" {
		query += "password = $" + fmt.Sprint(count) + ", "
		args = append(args, dtoUpdateUser.Password)
		count++
	}

	// Remove the trailing comma and add the WHERE condition
	if len(args) == 0 {
		logger.Warning(fmt.Sprintf("No data to update for user: %s", dtoUpdateUser.Username))
		return 0, exception.NewApiException(400, "No data to update")
	}

	// Using query[:len(query)-2], we remove the last comma and space, then add the WHERE condition
	query = query[:len(query)-2] + " WHERE username = $" + fmt.Sprint(count)
	args = append(args, dtoUpdateUser.Username)

	result, err := u.db.Exec(query, args...)
	if err != nil {
		logger.Error(fmt.Sprintf("Error updating user %s: %s", dtoUpdateUser.Username, err.Error()))
		return 0, exception.NewApiException(500, "Error updating user")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting affected rows for user %s: %s", dtoUpdateUser.Username, err.Error()))
		return 0, exception.NewApiException(500, "Error getting affected rows")
	}

	if rowsAffected == 0 {
		logger.Warning(fmt.Sprintf("User not found for update: %s", dtoUpdateUser.Username))
		return 0, exception.NewApiException(404, "User not found for update")
	}

	logger.Info(fmt.Sprintf("User successfully updated: %s, affected rows: %d", dtoUpdateUser.Username, rowsAffected))

	return rowsAffected, nil
}

func (u *UserPostgreSQLRepository) Delete(dtoDeleteUser *userDTO.UserDTO) (int64, *exception.ApiException) {
	logger.Info(fmt.Sprintf("Attempting to delete user: %s", dtoDeleteUser.Username))
	// Check if the user exists and that the password is correct before deleting the user
	_, errFind := u.find(dtoDeleteUser.Username, dtoDeleteUser.Password)
	if errFind != nil {
		logger.Warning(fmt.Sprintf("Error searching for user to delete %s: %s", dtoDeleteUser.Username, errFind.Message))
		return 0, errFind
	}

	query := "DELETE FROM users WHERE username = $1"
	result, err := u.db.Exec(query, dtoDeleteUser.Username)
	if err != nil {
		logger.Error(fmt.Sprintf("Error deleting user %s: %s", dtoDeleteUser.Username, err.Error()))
		return 0, exception.NewApiException(500, "Error deleting user")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting affected rows when deleting %s: %s", dtoDeleteUser.Username, err.Error()))
		return 0, exception.NewApiException(500, "Error getting affected rows")
	}

	if rowsAffected == 0 {
		logger.Warning(fmt.Sprintf("User not found for deletion: %s", dtoDeleteUser.Username))
		return 0, exception.NewApiException(404, "User not found for deletion")
	}

	logger.Info(fmt.Sprintf("User successfully deleted: %s", dtoDeleteUser.Username))
	return rowsAffected, nil
}

func (r *UserPostgreSQLRepository) checkUserIsCreated(dto *userDTO.UserDTO) *exception.ApiException {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)"
	var exists bool
	err := r.db.QueryRow(query, dto.Username, dto.Email).Scan(&exists)

	if err != nil {
		return exception.NewApiException(500, "Error verifying user existence")
	}

	if exists {
		return exception.NewApiException(400, "The username or email is already registered")
	}

	return nil
}
