package infrastructure

import (
	"api-upload-photos/src/domain"
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const (
	DBUser     = "admin"
	DBPassword = "password123"
	DBName     = "discorddb"
	DBPort     = "5432"
)

func AddImageToDataBase(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Error al obtener la imagen del formulario")
	}

	fileBytes, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error al abrir el archivo de imagen")
	}
	defer fileBytes.Close()

	fileData, err := io.ReadAll(fileBytes)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error al leer el archivo de imagen")
	}

	id := uuid.New()

	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable",
		DBUser, DBPassword, DBName, DBPort)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := insertImage(db, id.String(), fileData); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error al insertar la imagen en la base de datos")
	}

	response := domain.Response{
		Status: "success",
		Id:     id.String(),
	}

	return c.JSON(response)
}

func insertImage(db *sql.DB, id string, data []byte) error {
	_, err := db.Exec("INSERT INTO images (id, data) VALUES ($1, $2)", id, data)
	if err != nil {
		return err
	}
	return nil
}

func GetImageByID(c *fiber.Ctx) error {
	id := c.Params("id")

	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable",
		DBUser, DBPassword, DBName, DBPort)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	data := []byte{}
	err = db.QueryRow("SELECT data FROM images WHERE id = $1", id).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "Imagen no encontrada")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error al recuperar la imagen")
	}

	c.Set("Content-Type", "image/jpeg")

	_, err = io.Copy(c.Response().BodyWriter(), bytes.NewReader(data))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error al enviar la imagen como respuesta")
	}

	return nil
}