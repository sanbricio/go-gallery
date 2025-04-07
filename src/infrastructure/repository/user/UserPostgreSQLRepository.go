package user_repository

import (
	"database/sql"
	"fmt"
	"go-gallery/src/commons/exception"
	entity "go-gallery/src/domain/entities"
	"go-gallery/src/domain/entities/builder"
	"go-gallery/src/infrastructure/dto"

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

func (u *UserPostgreSQLRepository) Find(dtoLoginRequest *dto.DTOLoginRequest) (*dto.DTOUser, *exception.ApiException) {
	user, err := u.find(dtoLoginRequest.Username, dtoLoginRequest.Password)
	if err != nil {
		return nil, err
	}

	return dto.FromUser(user), nil
}

func (u *UserPostgreSQLRepository) find(username, password string) (*entity.User, *exception.ApiException) {
	query := "SELECT username, email, firstname,lastname, password FROM users WHERE username = $1"
	row := u.db.QueryRow(query, username)

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

	if err := user.CheckPasswordIntegrity(password); err != nil {
		return nil, exception.NewApiException(400, "Contraseña incorrecta")
	}

	return user, nil
}

func (u *UserPostgreSQLRepository) FindAndCheckJWT(claims *dto.DTOClaimsJwt) (*dto.DTOUser, *exception.ApiException) {
	query := "SELECT username, email FROM users WHERE username = $1"
	row := u.db.QueryRow(query, claims.Username)

	userDTO := new(dto.DTOUser)
	if err := row.Scan(&userDTO.Username, &userDTO.Email); err != nil {
		return nil, exception.NewApiException(404, "Usuario no encontrado")
	}

	if userDTO.Username != claims.Username || userDTO.Email != claims.Email {
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

func (u *UserPostgreSQLRepository) Update(dtoUpdateUser *dto.DTOUser) (int64, *exception.ApiException) {
	query := "UPDATE users SET "
	args := []any{}
	count := 1
	// Construimos la query y los argumentos para la actualización de manera dinámica
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

	// Eliminar la última coma y agregar condición WHERE
	if len(args) == 0 {
		return 0, exception.NewApiException(400, "No hay datos para actualizar")
	}

	// Con el query[:len(query)-2] eliminamos la última coma y espacio, despues de eso añadimos la condición WHERE
	query = query[:len(query)-2] + " WHERE username = $" + fmt.Sprint(count)
	args = append(args, dtoUpdateUser.Username)

	result, err := u.db.Exec(query, args...)
	if err != nil {
		return 0, exception.NewApiException(500, "Error al intentar actualizar el usuario")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, exception.NewApiException(500, "Error al intentar obtener el número de filas afectadas")
	}

	if rowsAffected == 0 {
		return 0, exception.NewApiException(404, "No se ha encontrado usuario para actualizar")
	}

	return rowsAffected, nil
}

func (u *UserPostgreSQLRepository) Delete(dtoDeleteUser *dto.DTOUser) (int64, *exception.ApiException) {
	// Comprobamos que el usuario exista y que la contraseña sea correcta para eliminar el usuario
	_, errFind := u.find(dtoDeleteUser.Username, dtoDeleteUser.Password)
	if errFind != nil {
		return 0, errFind
	}

	query := "DELETE FROM users WHERE username = $1"
	result, err := u.db.Exec(query, dtoDeleteUser.Username)
	if err != nil {
		return 0, exception.NewApiException(500, "Error al eliminar el usuario")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, exception.NewApiException(500, "Error al intentar obtener el número de filas afectadas")
	}

	if rowsAffected == 0 {
		return 0, exception.NewApiException(404, "No se ha encontrado usuario para eliminar")
	}

	return rowsAffected, nil
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
