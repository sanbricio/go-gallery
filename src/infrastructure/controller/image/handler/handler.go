package imageHandler

import (
	"go-gallery/src/commons/constants"
	"go-gallery/src/commons/exception"
	utilsImage "go-gallery/src/commons/utils/image"
	imageDTO "go-gallery/src/infrastructure/dto/image"
	"io"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strings"
)

func ProcessImageFile(fileInput *multipart.FileHeader, owner string) (*imageDTO.ImageUploadRequestDTO, *exception.ApiException) {

	fileExtension := filepath.Ext(fileInput.Filename)
	fileName := strings.TrimSuffix(fileInput.Filename, fileExtension)

	if !isValidExtension(fileExtension) {
		return nil, exception.NewApiException(400, "Formato de archivo no soportado. Solo se aceptan imágenes jpg, jpeg, png y webp")
	}

	rawData, err := encodeToRawBytes(fileInput)
	if err != nil {
		return nil, err
	}

	fileSizeHumanReadable := utilsImage.HumanizeBytes(uint64(fileInput.Size))

	return &imageDTO.ImageUploadRequestDTO{
		Name:           fileName,
		Extension:      fileExtension,
		Size:           fileSizeHumanReadable,
		RawContentFile: rawData,
		Owner:          owner,
	}, nil
}

func isValidExtension(extension string) bool {
	validExtensions := []string{constants.JPG_EXTENSION, constants.JPEG_EXTENSION, constants.PNG_EXTENSION, constants.WEBP_EXTENSION}
	return slices.Contains(validExtensions, extension)
}

func encodeToRawBytes(fileInput *multipart.FileHeader) ([]byte, *exception.ApiException) {
	fileBytes, err := fileInput.Open()
	if err != nil {
		return nil, exception.NewApiException(500, "Error al abrir el archivo de imagen")
	}

	defer fileBytes.Close()

	return readAllFile(fileBytes)
}


func readAllFile(file multipart.File) ([]byte, *exception.ApiException) {
    fileData, err := io.ReadAll(file)
    if err != nil {
        return nil, exception.NewApiException(500, "Error al leer el archivo de imagen")
    }
    return fileData, nil
}