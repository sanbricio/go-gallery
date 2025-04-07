package handler

import (
	"encoding/base64"
	"go-gallery/src/commons/exception"
	"go-gallery/src/infrastructure/dto"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"slices"

	"github.com/dustin/go-humanize"
)

func ProcessImageFile(fileInput *multipart.FileHeader, owner string) (*dto.DTOImage, *exception.ApiException) {

	fileExtension := filepath.Ext(fileInput.Filename)
	fileName := strings.TrimSuffix(fileInput.Filename, fileExtension)

	if !isValidExtension(fileExtension) {
		return nil, exception.NewApiException(400, "Formato de archivo no soportado. Solo se aceptan imágenes jpg, jpeg, png y webp")
	}

	encoded, err := encodeToBase64(fileInput)
	if err != nil {
		return nil, err
	}

	fileSizeHumanReadable := humanize.Bytes(uint64(fileInput.Size))

	processedFile := &dto.DTOImage{
		Name:        fileName,
		Extension:   fileExtension,
		Size:        fileSizeHumanReadable,
		ContentFile: encoded,
		Owner:       owner,
	}

	return processedFile, nil
}

func isValidExtension(extension string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".webp"}

	return slices.Contains(validExtensions, extension)
}

func encodeToBase64(fileInput *multipart.FileHeader) (string, *exception.ApiException) {
	fileBytes, err := fileInput.Open()
	if err != nil {
		return "", exception.NewApiException(500, "Error al abrir el archivo de imagen")
	}

	defer fileBytes.Close()

	fileData, err := io.ReadAll(fileBytes)
	if err != nil {
		return "", exception.NewApiException(500, "Error al leer el archivo de imagen")
	}

	return base64.StdEncoding.EncodeToString(fileData), nil
}
