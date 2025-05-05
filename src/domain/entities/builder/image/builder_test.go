package imageBuilder

import (
	"fmt"
	"testing"

	imageEntity "go-gallery/src/domain/entities/image"
	imageDTO "go-gallery/src/infrastructure/dto/image"

	"github.com/stretchr/testify/assert"
)

const (
	UNEXPECTED_ERROR string = "An unexpected error occurred while building the image: %v"
)

var baseDTO *imageDTO.ImageDTO

// Initialize the mock DTO
func init() {
	id := "valid-id"
	baseDTO = &imageDTO.ImageDTO{
		Id:          &id,
		ThumbnailId: "valid-thumbnail-id",
		Name:        "valid-name",
		Extension:   "jpg",
		ContentFile: "base64content",
		Owner:       "valid-owner",
		Size:        "1024",
	}
}

func TestImageBuilderEmptyFields(t *testing.T) {
	// Case: 'id' is empty
	dto := copyDTO()
	dto.Id = nil
	assertBuilderException(t, dto, "id")

	// Case: 'thumbnailId' is empty
	dto = copyDTO()
	dto.ThumbnailId = ""
	assertBuilderException(t, dto, "thumbnailId")

	// Case: 'name' is empty
	dto = copyDTO()
	dto.Name = ""
	assertBuilderException(t, dto, "name")

	// Case: 'extension' is empty
	dto = copyDTO()
	dto.Extension = ""
	assertBuilderException(t, dto, "extension")

	// Case: 'contentFile' is empty
	dto = copyDTO()
	dto.ContentFile = ""
	assertBuilderException(t, dto, "contentFile")

	// Case: 'owner' is empty
	dto = copyDTO()
	dto.Owner = ""
	assertBuilderException(t, dto, "owner")

	// Case: 'size' is empty
	dto = copyDTO()
	dto.Size = ""
	assertBuilderException(t, dto, "size")
}

func TestImageBuilderNew(t *testing.T) {
	image, err := NewImageBuilder().
		SetName(baseDTO.Name).
		SetExtension(baseDTO.Extension).
		SetOwner(baseDTO.Owner).
		SetSize(baseDTO.Size).
		SetContentFile(baseDTO.ContentFile).
		BuildNew()

	assert.Nil(t, err, fmt.Sprintf(UNEXPECTED_ERROR, err), err)
	compareCommonFieldsImages(t, baseDTO, image)

	// Error when building without a specific field
	_, err = NewImageBuilder().
		SetName(baseDTO.Name).
		SetExtension(baseDTO.Extension).
		SetOwner(baseDTO.Owner).
		SetSize(baseDTO.Size).
		SetContentFile(baseDTO.ContentFile).
		BuildNew()

	assert.Error(t, err, "An error was expected when building the image due to missing 'thumbnailID'")
}

func TestImageBuilderFromImageUploadRequestDTO(t *testing.T) {
	dto := &imageDTO.ImageUploadRequestDTO{
		Name:           baseDTO.Name,
		Extension:      baseDTO.Extension,
		RawContentFile: []byte(baseDTO.ContentFile),
		Owner:          baseDTO.Owner,
		Size:           baseDTO.Size,
	}

	image, err := NewImageBuilder().FromImageUploadRequestDTO(dto).
		SetContentFile(baseDTO.ContentFile).
		BuildNew()

	assert.Nil(t, err, fmt.Sprintf(UNEXPECTED_ERROR, err), err)
	compareCommonFieldsImages(t, baseDTO, image)
}

func TestImageBuilderWithSetValues(t *testing.T) {
	image, err := NewImageBuilder().
		SetId(baseDTO.Id).
		SetName(baseDTO.Name).
		SetExtension(baseDTO.Extension).
		SetOwner(baseDTO.Owner).
		SetSize(baseDTO.Size).
		SetContentFile(baseDTO.ContentFile).
		Build()

	assert.Nil(t, err, fmt.Sprintf(UNEXPECTED_ERROR, err), err)
	compareAllFieldsImages(t, baseDTO, image)
}

func compareAllFieldsImages(t *testing.T, expected *imageDTO.ImageDTO, actual *imageEntity.Image) {
	assert.Equal(t, expected.Id, actual.GetId(), "expected id does not match")
	compareCommonFieldsImages(t, expected, actual)
}

func compareCommonFieldsImages(t *testing.T, expected *imageDTO.ImageDTO, actual *imageEntity.Image) {
	assert.Equal(t, expected.Name, actual.GetName(), "expected name does not match")
	assert.Equal(t, expected.Extension, actual.GetExtension(), "expected extension does not match")
	assert.Equal(t, expected.Owner, actual.GetOwner(), "expected owner does not match")
	assert.Equal(t, expected.Size, actual.GetSize(), "expected size does not match")
	assert.Equal(t, expected.ContentFile, actual.GetContentFile(), "expected contentFile does not match")
}

func assertBuilderException(t *testing.T, dto *imageDTO.ImageDTO, field string) {
	_, err := NewImageBuilder().
		FromDTO(dto).
		Build()

	assert.Error(t, err, "An error was expected when trying to create an Image with the field '%v' empty", field)
	if err != nil {
		assert.Equal(t, field, err.Field, "An error for the field '%v' was expected, but got: %v", field, err)
	}
}

// Function to copy the base DTO for each test case
func copyDTO() *imageDTO.ImageDTO {
	return &imageDTO.ImageDTO{
		Id:          baseDTO.Id,
		ThumbnailId: baseDTO.ThumbnailId,
		Name:        baseDTO.Name,
		Extension:   baseDTO.Extension,
		ContentFile: baseDTO.ContentFile,
		Owner:       baseDTO.Owner,
		Size:        baseDTO.Size,
	}
}
