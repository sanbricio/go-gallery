package imageHandler

import (
	"encoding/base64"
	"go-gallery/src/commons/exception"
	imageDTO "go-gallery/src/infrastructure/dto/image"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"slices"

	"github.com/dustin/go-humanize"
)

func ProcessImageFile(fileInput *multipart.FileHeader, owner string) (*imageDTO.ImageUploadRequestDTO, *exception.ApiException) {

	fileExtension := filepath.Ext(fileInput.Filename)
	fileName := strings.TrimSuffix(fileInput.Filename, fileExtension)

	if !isValidExtension(fileExtension) {
		return nil, exception.NewApiException(400, "Formato de archivo no soportado. Solo se aceptan imágenes jpg, jpeg, png y webp")
	}

	encoded, rawData, err := encodeToBase64(fileInput)
	if err != nil {
		return nil, err
	}

	fileSizeHumanReadable := humanize.Bytes(uint64(fileInput.Size))

	return &imageDTO.ImageUploadRequestDTO{
		Name:           fileName,
		Extension:      fileExtension,
		Size:           fileSizeHumanReadable,
		ContentFile:    encoded,
		RawContentFile: rawData,
		Owner:          owner,
	}, nil
}

func isValidExtension(extension string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".webp"}

	return slices.Contains(validExtensions, extension)
}

func encodeToBase64(fileInput *multipart.FileHeader) (string, []byte, *exception.ApiException) {
	fileBytes, err := fileInput.Open()
	if err != nil {
		return "", nil, exception.NewApiException(500, "Error al abrir el archivo de imagen")
	}

	defer fileBytes.Close()

	fileData, err := io.ReadAll(fileBytes)
	if err != nil {
		return "", nil, exception.NewApiException(500, "Error al leer el archivo de imagen")
	}

	return base64.StdEncoding.EncodeToString(fileData), fileData, nil
}
