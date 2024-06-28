package handler

import (
	"api-upload-photos/src/commons/exception"
	"encoding/base64"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
)

// TODO Pensar si esta solución es correcta (dominio) 
type ProcessedImage struct {
	FileName              string
	FileExtension         string
	FileSizeHumanReadable string
	EncodedData           string
}

func ProcessImageFile(fileInput *multipart.FileHeader) (*ProcessedImage, *exception.ApiException) {

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

	processedFile := &ProcessedImage{
		FileName:              fileName,
		FileExtension:         fileExtension,
		FileSizeHumanReadable: fileSizeHumanReadable,
		EncodedData:           encoded,
	}

	return processedFile, nil
}

func isValidExtension(extension string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".webp"}

	for _, validExt := range validExtensions {
		if extension == validExt {
			return true
		}
	}
	return false
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
