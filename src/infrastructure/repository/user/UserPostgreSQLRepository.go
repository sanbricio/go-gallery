package user_repository

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/domain/entities/builder"
	"api-upload-photos/src/infrastructure/dto"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const UserPostgreSQLRepositoryKey = "UserPostgreSQLRepository"

type UserPostgreSQLRepository struct {
	db *sql.DB
}

func NewUserPostgreSQLRepository(args map[string]string) UserRepository {
	user := args["POSTGRESQL_USER"]
	password := args["POSTGRESQL_PASSWORD"]
	dbName := args["POSTGRESQL_DB"]
	host := args["POSTGRESQL_HOST"]
	port := args["POSTGRESQL_PORT"]

	// Construimos la URL de conexión
	urlConnection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

	db, err := sql.Open("postgres", urlConnection)
	if err != nil {
		panic(fmt.Sprintf("No se ha podido conectar a PostgreSQL: %s", err.Error()))
	}

	return &UserPostgreSQLRepository{db: db}
}

func (u *UserPostgreSQLRepository) Find(dtoLoginRequest *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	query := "SELECT username, email, firstname,lastname, password FROM users WHERE username = $1"
	row := u.db.QueryRow(query, dtoLoginRequest.Username)

	userDTO := new(dto.DTOUser)
	if err := row.Scan(&userDTO.Username, &userDTO.Email, &userDTO.Firstname, &userDTO.Lastname, &userDTO.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewApiException(404, "Usuario no encontrado")
		}
		return nil, exception.NewApiException(500, "Error al buscar usuario")
	}

	user, errBuilder := builder.NewUserBuilder().
		FromDTO(userDTO).
		Build()

	if errBuilder != nil {
		return nil, exception.NewApiException(500, errBuilder.Error())
	}

	if err := user.CheckPasswordIntegrity(dtoLoginRequest.Password); err != nil {
		return nil, exception.NewApiException(404, "Contraseña incorrecta")
	}

	return dto.FromUser(user), nil
}

func (u *UserPostgreSQLRepository) FindAndCheckJWT(claims *dto.DTOClaimsJwt) (*dto.DTOUser, *exception.ApiException) {
	query := "SELECT username, email, firstname FROM users WHERE username = $1"
	row := u.db.QueryRow(query, claims.Username)

	userDTO := new(dto.DTOUser)
	if err := row.Scan(&userDTO.Username, &userDTO.Email, &userDTO.Firstname); err != nil {
		return nil, exception.NewApiException(404, "Usuario no encontrado")
	}

	if userDTO.Username != claims.Username || userDTO.Email != claims.Email || userDTO.Firstname != claims.Firstname {
		return nil, exception.NewApiException(403, "Los datos proporcionados no coinciden con el usuario autenticado")
	}

	return userDTO, nil
}

func (u *UserPostgreSQLRepository) Insert(dtoRegisterRequest *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	if err := u.checkUserIsCreated(dtoRegisterRequest); err != nil {
		return nil, err
	}

	user, err := builder.NewUserBuilder().
		FromDTO(dtoRegisterRequest).
		Build()

	if err != nil {
		return nil, exception.NewApiException(500, err.Error())
	}

	query := "INSERT INTO users (username, email, firstname, lastname, password) VALUES ($1, $2, $3, $4, $5)"
	_, errDb := u.db.Exec(query, user.GetUsername(), user.GetEmail(), user.GetFirstname(), user.GetLastname(), user.GetPassword())
	if errDb != nil {
		return nil, exception.NewApiException(500, "Error al insertar el usuario")
	}

	return dto.FromUser(user), nil
}

func (r *UserPostgreSQLRepository) Update(dtoUpdateUser *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	query := "UPDATE users SET email = $1, firstname = $2, lastname = $3  password = $4 WHERE username = $5"
	_, err := r.db.Exec(query, dtoUpdateUser.Email, dtoUpdateUser.Firstname, dtoUpdateUser.Lastname, dtoUpdateUser.Password, dtoUpdateUser.Username)
	if err != nil {
		return nil, exception.NewApiException(500, "Error al actualizar el usuario")
	}

	return dtoUpdateUser, nil
}

func (r *UserPostgreSQLRepository) Delete(dtoDeleteUser *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	query := "DELETE FROM users WHERE username = $1"
	_, err := r.db.Exec(query, dtoDeleteUser.Username)
	if err != nil {
		return nil, exception.NewApiException(500, "Error al eliminar el usuario")
	}

	return dtoDeleteUser, nil
}

func (r *UserPostgreSQLRepository) checkUserIsCreated(dto *dto.DTOUser) *exception.ApiException {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)"
	var exists bool
	err := r.db.QueryRow(query, dto.Username, dto.Email).Scan(&exists)

	if err != nil {
		return exception.NewApiException(500, "Error al verificar existencia del usuario")
	}

	if exists {
		return exception.NewApiException(400, "El usuario o correo ya están registrados")
	}

	return nil
}
