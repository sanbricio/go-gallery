package infrastructure

import (
	"api-upload-photos/src/domain"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type RepositoryMemory struct {
}


func (r *RepositoryMemory) Find(id string) (*domain.DTOImage, error) {
	files, err := os.ReadDir("data")

	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Error al leer el directorio")
	}
	var dtoImage *domain.DTOImage
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			idFilename := strings.TrimSuffix(file.Name(), ".json")

			if id == idFilename {
				data, err := os.ReadFile(filepath.Join("data", file.Name()))
				if err != nil {
					return nil, fiber.NewError(fiber.StatusInternalServerError, "Error al leer el archivo JSON")
				}

				err = json.Unmarshal(data, &dtoImage)
				if err != nil {
					return nil, fiber.NewError(fiber.StatusInternalServerError, "Error al parsear el archivo JSON")
				}

				return dtoImage, nil

			}
		}
	}
	return nil, fiber.NewError(fiber.StatusNotFound, "Imagen no encontrada")
}

func (r *RepositoryMemory) Insert(fileInput *multipart.FileHeader) (*domain.Response, error) {
	fileCompleteName := strings.Split(fileInput.Filename, ".")

	name := fileCompleteName[0]

	extension := fileCompleteName[1]

	fileBytes, err := fileInput.Open()
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Error al abrir el archivo de imagen")
	}
	defer fileBytes.Close()

	fileData, err := io.ReadAll(fileBytes)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Error al leer el archivo de imagen")
	}

	encoded := base64.StdEncoding.EncodeToString(fileData)

	fileSizeHumanReadable := humanize.Bytes(uint64(fileInput.Size))

	image := domain.NewImage(uuid.New().String(), name, extension, encoded, "SANTI", fileSizeHumanReadable)

	dto := image.ToDto(image)

	err = persist(dto)

	if err != nil {
		return nil, err
	}

	response := &domain.Response{
		Status: "success",
		Id:     dto.Id,
	}

	return response, nil
}

func persist(image *domain.DTOImage) error {

	err := os.MkdirAll("data", 0755)

	if err != nil {
		return err
	}

	filename := fmt.Sprintf("data/%s.json", image.Id)

	data, err := json.Marshal(image)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)

	if err != nil {
		return err
	}

	return nil
}
