package imageHandler

import (
	"go-gallery/src/commons/constants"
	"go-gallery/src/commons/exception"
	utilsImage "go-gallery/src/commons/utils/image"
	imageDTO "go-gallery/src/infrastructure/dto/image"
	"go-gallery/src/infrastructure/logger"
	"io"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

func ProcessImageFile(fileInput *multipart.FileHeader, owner string) (*imageDTO.ImageUploadRequestDTO, *exception.ApiException) {
	logger.Instance().Info("Starting image file processing: filename=" + fileInput.Filename + ", owner=" + owner)

	fileExtension := filepath.Ext(fileInput.Filename)
	fileName := strings.TrimSuffix(fileInput.Filename, fileExtension)
	logger.Instance().Info("Extracted filename and extension: name=" + fileName + ", extension=" + fileExtension)

	if !isValidExtension(fileExtension) {
		logger.Instance().Warning("Invalid file extension detected: extension=" + fileExtension)
		return nil, exception.NewApiException(400, "Unsupported file format. Only jpg, jpeg, png, and webp images are accepted.")
	}

	rawData, err := encodeToRawBytes(fileInput)
	if err != nil {
		logger.Instance().Error("Failed to convert file to bytes: error=" + err.Message)
		return nil, err
	}

	fileSizeHumanReadable := utilsImage.HumanizeBytes(uint64(fileInput.Size))
	logger.Instance().Info("File processed successfully: name=" + fileName + ", extension=" + fileExtension + ", size=" + fileSizeHumanReadable)

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
	isValid := slices.Contains(validExtensions, extension)
	logger.Instance().Info("Checking file extension: extension=" + extension + ", isValid=" + strconv.FormatBool(isValid))
	return isValid
}

func encodeToRawBytes(fileInput *multipart.FileHeader) ([]byte, *exception.ApiException) {
	logger.Instance().Info("Opening file for byte conversion: filename=" + fileInput.Filename)

	fileBytes, err := fileInput.Open()
	if err != nil {
		logger.Instance().Error("Failed to open file: filename=" + fileInput.Filename + ", error=" + err.Error())
		return nil, exception.NewApiException(500, "Error opening the image file")
	}
	defer fileBytes.Close()

	return readAllFile(fileBytes)
}

func readAllFile(file multipart.File) ([]byte, *exception.ApiException) {
	logger.Instance().Info("Reading entire file into memory")
	fileData, err := io.ReadAll(file)
	if err != nil {
		logger.Instance().Error("Failed to read file content")
		return nil, exception.NewApiException(500, "Error reading the image file")
	}
	logger.Instance().Info("File read successfully: sizeBytes=" + strconv.Itoa(len(fileData)))
	return fileData, nil
}
