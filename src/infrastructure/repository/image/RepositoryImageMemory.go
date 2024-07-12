package repository

import (
	"api-upload-photos/src/commons/exception"
	entity "api-upload-photos/src/domain/entities"
	"api-upload-photos/src/infrastructure/dto"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type RepositoryImageMemory struct {
}

func (r *RepositoryImageMemory) Find(id string) (*dto.DTOImage, *exception.ApiException) {
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

func (r *RepositoryImageMemory) Insert(processedImage *dto.DTOProcessedImage) (*dto.DTOImage, *exception.ApiException) {

	image := entity.NewImage(processedImage.FileName, processedImage.FileExtension, processedImage.EncodedData, "SANTI", processedImage.FileSizeHumanReadable)

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

func (r *RepositoryImageMemory) Delete(id string) (*dto.DTOImage, *exception.ApiException) {
	image, err := r.Find(id)
	if err != nil {
		return nil, err
	}

	filename := filepath.Join("data", fmt.Sprintf("%s.json", id))

	errRemove := os.Remove(filename)
	if errRemove != nil {
		return nil, exception.NewApiException(500, fmt.Sprintf("Error al eliminar el archivo de imagen: %s", filename))
	}

	return image, nil
}
