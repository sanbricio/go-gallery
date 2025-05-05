package thumbnailImageBuilder

import (
	"testing"

	thumbnailImageEntity "go-gallery/src/domain/entities/image/thumbnailImage"
	imageDTO "go-gallery/src/infrastructure/dto/image"
	thumbnailImageDTO "go-gallery/src/infrastructure/dto/image/thumbnailImage"

	"github.com/stretchr/testify/assert"
)

const (
	UNEXPECTED_ERROR string = "No se esperaba un error al construir la miniatura: %v"
)

var baseThumbnailDTO *thumbnailImageDTO.ThumbnailImageDTO

// Inicializamos el mock dto
func init() {
	id := "valid-id"
	baseThumbnailDTO = &thumbnailImageDTO.ThumbnailImageDTO{
		Id:          &id,
		ImageID:     &id,
		Name:        "valid-name",
		Extension:   "jpg",
		ContentFile: "base64content",
		Owner:       "valid-owner",
		Size:        "1024",
	}
}

func TestThumbnailImageBuilderEmptyFields(t *testing.T) {
	// Caso: 'id' vacío (solo aplica para Build)
	dto := copyThumbnailDTO()
	dto.Id = nil
	assertThumbnailBuilderException(t, dto, "id")

	// Caso: 'name' vacío
	dto = copyThumbnailDTO()
	dto.Name = ""
	assertThumbnailBuilderException(t, dto, "name")

	// Caso: 'extension' vacío
	dto = copyThumbnailDTO()
	dto.Extension = ""
	assertThumbnailBuilderException(t, dto, "extension")

	// Caso: 'contentFile' vacío
	dto = copyThumbnailDTO()
	dto.ContentFile = ""
	assertThumbnailBuilderException(t, dto, "contentFile")

	// Caso: 'owner' vacío
	dto = copyThumbnailDTO()
	dto.Owner = ""
	assertThumbnailBuilderException(t, dto, "owner")

	// Caso: 'size' vacío
	dto = copyThumbnailDTO()
	dto.Size = ""
	assertThumbnailBuilderException(t, dto, "size")
}

func TestThumbnailImageBuilderNew(t *testing.T) {
	image, err := NewThumbnailImageBuilder().
		SetImageID(baseThumbnailDTO.ImageID).
		SetName(baseThumbnailDTO.Name).
		SetExtension(baseThumbnailDTO.Extension).
		SetOwner(baseThumbnailDTO.Owner).
		SetSize(baseThumbnailDTO.Size).
		SetContentFile(baseThumbnailDTO.ContentFile).
		BuildNew()

	assert.Error(t, err, UNEXPECTED_ERROR, err)
	compareCommonFieldsThumbnailImage(t, baseThumbnailDTO, image)

	// Error: falta 'name'
	_, err = NewThumbnailImageBuilder().
		SetExtension(baseThumbnailDTO.Extension).
		SetOwner(baseThumbnailDTO.Owner).
		SetSize(baseThumbnailDTO.Size).
		SetContentFile(baseThumbnailDTO.ContentFile).
		BuildNew()

	assert.Error(t, err, "Se esperaba un error al construir la miniatura debido a que faltaba 'name'")
}

func TestThumbnailImageBuilderFromImageUploadRequestDTO(t *testing.T) {
	dto := &imageDTO.ImageUploadRequestDTO{
		Name:      baseThumbnailDTO.Name,
		Extension: baseThumbnailDTO.Extension,
		Owner:     baseThumbnailDTO.Owner,
		Size:      baseThumbnailDTO.Size,
	}

	image, err := NewThumbnailImageBuilder().
		SetImageID(baseThumbnailDTO.ImageID).
		SetContentFile(baseThumbnailDTO.ContentFile).
		FromImageUploadRequestDTO(dto).
		BuildNew()

	assert.Error(t, err, UNEXPECTED_ERROR, err)
	compareCommonFieldsThumbnailImage(t, baseThumbnailDTO, image)
}

func TestThumbnailImageBuilderWithSetValues(t *testing.T) {
	image, err := NewThumbnailImageBuilder().
		SetId(baseThumbnailDTO.Id).
		SetImageID(baseThumbnailDTO.ImageID).
		SetName(baseThumbnailDTO.Name).
		SetExtension(baseThumbnailDTO.Extension).
		SetOwner(baseThumbnailDTO.Owner).
		SetSize(baseThumbnailDTO.Size).
		SetContentFile(baseThumbnailDTO.ContentFile).
		Build()

	assert.Error(t, err, UNEXPECTED_ERROR, err)
	compareAllFieldsThumbnailImage(t, baseThumbnailDTO, image)
}

func compareAllFieldsThumbnailImage(t *testing.T, expected *thumbnailImageDTO.ThumbnailImageDTO, actual *thumbnailImageEntity.ThumbnailImage) {
	if expected.Id == nil {
		assert.Nil(t, actual.GetId(), "expected id nil, but got %v", actual.GetId())
	} else {
		assert.NotNil(t, actual.GetId(), "expected id %v, but got nil", *expected.Id)
		assert.Equal(t, *expected.Id, *actual.GetId(), "expected id %v, but got %v", *expected.Id, *actual.GetId())
	}

	compareCommonFieldsThumbnailImage(t, expected, actual)
}

func compareCommonFieldsThumbnailImage(t *testing.T, expected *thumbnailImageDTO.ThumbnailImageDTO, actual *thumbnailImageEntity.ThumbnailImage) {
	assert.Equal(t, expected.Name, actual.GetName(), "expected name %v, but got %v", expected.Name, actual.GetName())
	assert.Equal(t, expected.Extension, actual.GetExtension(), "expected extension %v, but got %v", expected.Extension, actual.GetExtension())
	assert.Equal(t, expected.Owner, actual.GetOwner(), "expected owner %v, but got %v", expected.Owner, actual.GetOwner())
	assert.Equal(t, expected.Size, actual.GetSize(), "expected size %v, but got %v", expected.Size, actual.GetSize())
	assert.Equal(t, expected.ContentFile, actual.GetContentFile(), "expected contentFile %v, but got %v", expected.ContentFile, actual.GetContentFile())
}

func assertThumbnailBuilderException(t *testing.T, dto *thumbnailImageDTO.ThumbnailImageDTO, field string) {
	_, err := NewThumbnailImageBuilder().
		FromDTO(dto).
		Build()

	assert.Error(t, err, "Se esperaba un error al intentar crear una ThumbnailImage con el campo '%v' vacío", field)
	assert.Equal(t, field, err.Field, "Se esperaba un error del campo '%v', pero se obtuvo: %v", field, err.Field)
}

// Función para copiar el DTO base para cada caso de prueba
func copyThumbnailDTO() *thumbnailImageDTO.ThumbnailImageDTO {
	return &thumbnailImageDTO.ThumbnailImageDTO{
		Id:          baseThumbnailDTO.Id,
		ImageID:     baseThumbnailDTO.ImageID,
		Name:        baseThumbnailDTO.Name,
		Extension:   baseThumbnailDTO.Extension,
		ContentFile: baseThumbnailDTO.ContentFile,
		Owner:       baseThumbnailDTO.Owner,
		Size:        baseThumbnailDTO.Size,
	}
}
