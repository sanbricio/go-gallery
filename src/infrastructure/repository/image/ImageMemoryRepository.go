package image_repository

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/domain/entities/builder"
	"api-upload-photos/src/infrastructure/dto"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const ImageMemoryRepositoryKey = "ImageMemoryRepository"

type ImageMemoryRepository struct {
}

func NewImageMemoryRepository() ImageRepository {
	return new(ImageMemoryRepository)
}

func (r *ImageMemoryRepository) Find(dtoFind *dto.DTOImage) (*dto.DTOImage, *exception.ApiException) {
	files, err := os.ReadDir("data")
	if err != nil {
		return nil, exception.NewApiException(500, "Error al leer el directorio")
	}

	var dtoImage *dto.DTOImage
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			idFilename := strings.TrimSuffix(file.Name(), ".json")

			if dtoFind.IdImage == idFilename {
				data, err := os.ReadFile(filepath.Join("data", file.Name()))
				if err != nil {
					return nil, exception.NewApiException(500, "Error al leer el archivo JSON")
				}

				err = json.Unmarshal(data, &dtoImage)
				if err != nil {
					return nil, exception.NewApiException(500, "Error al parsear el archivo JSON")
				}

				return dtoImage, nil

			}
		}
	}
	return nil, exception.NewApiException(404, "Imagen no encontrada")
}

func (r *ImageMemoryRepository) find(id string) (*dto.DTOImage, *exception.ApiException) {
	files, err := os.ReadDir("data")
	if err != nil {
		return nil, exception.NewApiException(500, "Error al leer el directorio")
	}

	var dtoImage *dto.DTOImage
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			idFilename := strings.TrimSuffix(file.Name(), ".json")

			if id == idFilename {
				data, err := os.ReadFile(filepath.Join("data", file.Name()))
				if err != nil {
					return nil, exception.NewApiException(500, "Error al leer el archivo JSON")
				}

				err = json.Unmarshal(data, &dtoImage)
				if err != nil {
					return nil, exception.NewApiException(500, "Error al parsear el archivo JSON")
				}

				return dtoImage, nil

			}
		}
	}
	return nil, exception.NewApiException(404, "Imagen no encontrada")
}

func (r *ImageMemoryRepository) Insert(dtoInsertImage *dto.DTOImage) (*dto.DTOImage, *exception.ApiException) {

	image, errBuilder := builder.NewImageBuilder().
		FromDTO(dtoInsertImage).
		Build()

	if errBuilder != nil {
		errorMessage := fmt.Sprintf("Error al construir la imagen: %s", errBuilder.Error())
		return nil, exception.NewApiException(404, errorMessage)
	}

	dto := dto.FromImage(image)

	errPersist := persist(dto)
	if errPersist != nil {
		return nil, errPersist
	}

	return dto, nil
}

func persist(image *dto.DTOImage) *exception.ApiException {

	err := os.MkdirAll("data", 0755)
	if err != nil {
		return exception.NewApiException(500, "Error al crear el directorio de almacenamiento de datos")
	}

	filename := fmt.Sprintf("data/%s.json", image.IdImage)

	data, err := json.Marshal(image)
	if err != nil {
		return exception.NewApiException(500, "Error al convertir la imagen a JSON")
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return exception.NewApiException(500, "Error al escribir el archivo de imagen")
	}

	return nil
}

func (r *ImageMemoryRepository) Delete(dto *dto.DTOImage) (*dto.DTOImage, *exception.ApiException) {
	image, err := r.find(dto.IdImage)
	if err != nil {
		return nil, err
	}

	filename := filepath.Join("data", fmt.Sprintf("%s.json", dto.IdImage))

	errRemove := os.Remove(filename)
	if errRemove != nil {
		return nil, exception.NewApiException(500, fmt.Sprintf("Error al eliminar el archivo de imagen: %s", filename))
	}

	return image, nil
}
