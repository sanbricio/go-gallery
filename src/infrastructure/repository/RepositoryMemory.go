package infrastructure

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/domain/dto"
	entity "api-upload-photos/src/domain/entities"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type RepositoryMemory struct {
}

func (r *RepositoryMemory) Find(id string) (*dto.DTOImage, *exception.ApiException) {
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

func (r *RepositoryMemory) Insert(fileInput *multipart.FileHeader) (*entity.Response, *exception.ApiException) {
	//TODO Mirar bien el tratamiento de extensiones y el split de la manera correspondiente (no pueden pasar archivos .zip por ejemplo)
	fileCompleteName := strings.Split(fileInput.Filename, ".")
	name := fileCompleteName[0]
	extension := fileCompleteName[1]

	fileBytes, err := fileInput.Open()
	if err != nil {
		return nil, exception.NewApiException(500, "Error al abrir el archivo de imagen")
	}
	defer fileBytes.Close()

	fileData, err := io.ReadAll(fileBytes)
	if err != nil {
		return nil, exception.NewApiException(500, "Error al leer el archivo de imagen")
	}

	encoded := base64.StdEncoding.EncodeToString(fileData)

	fileSizeHumanReadable := humanize.Bytes(uint64(fileInput.Size))

	image := entity.NewImage(uuid.New().String(), name, extension, encoded, "SANTI", fileSizeHumanReadable)

	dto := dto.FromImage(image)

	errPersist := persist(dto)
	if errPersist != nil {
		return nil, errPersist
	}

	response := entity.NewResponse(dto.Id, "success")

	return response, nil
}

func persist(image *dto.DTOImage) *exception.ApiException {

	err := os.MkdirAll("data", 0755)

	if err != nil {
		return exception.NewApiException(500, "Error al crear el directorio de almacenamiento de datos")
	}

	filename := fmt.Sprintf("data/%s.json", image.Id)

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
