package thumbnailImageBuilder

import (
	thumbnailImageEntity "go-gallery/src/domain/entities/image/thumbnailImage"
	imageDTO "go-gallery/src/infrastructure/dto/image"
	thumbnailImageDTO "go-gallery/src/infrastructure/dto/image/thumbnailImage"
	"testing"
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
		SetName(baseThumbnailDTO.Name).
		SetExtension(baseThumbnailDTO.Extension).
		SetOwner(baseThumbnailDTO.Owner).
		SetSize(baseThumbnailDTO.Size).
		SetContentFile(baseThumbnailDTO.ContentFile).
		BuildNew()

	if err != nil {
		t.Errorf(UNEXPECTED_ERROR, err.Error())
		t.FailNow()
	}

	compareCommonFieldsThumbnailImage(t, baseThumbnailDTO, image)

	// Error: falta 'name'
	_, err = NewThumbnailImageBuilder().
		SetExtension(baseThumbnailDTO.Extension).
		SetOwner(baseThumbnailDTO.Owner).
		SetSize(baseThumbnailDTO.Size).
		SetContentFile(baseThumbnailDTO.ContentFile).
		BuildNew()

	if err == nil {
		t.Errorf("Se esperaba un error al construir la miniatura debido a que faltaba 'name'")
		t.FailNow()
	}
}

func TestThumbnailImageBuilderFromImageUploadRequestDTO(t *testing.T) {

	dto := &imageDTO.ImageUploadRequestDTO{
		Name:      baseThumbnailDTO.Name,
		Extension: baseThumbnailDTO.Extension,
		Owner:     baseThumbnailDTO.Owner,
		Size:      baseThumbnailDTO.Size,
	}

	image, err := NewThumbnailImageBuilder().
		SetContentFile(baseThumbnailDTO.ContentFile).
		FromImageUploadRequestDTO(dto).
		BuildNew()

	if err != nil {
		t.Errorf(UNEXPECTED_ERROR, err.Error())
		t.FailNow()
	}

	compareCommonFieldsThumbnailImage(t, baseThumbnailDTO, image)
}

func TestThumbnailImageBuilderWithSetValues(t *testing.T) {
	image, err := NewThumbnailImageBuilder().
		SetId(baseThumbnailDTO.Id).
		SetName(baseThumbnailDTO.Name).
		SetExtension(baseThumbnailDTO.Extension).
		SetOwner(baseThumbnailDTO.Owner).
		SetSize(baseThumbnailDTO.Size).
		SetContentFile(baseThumbnailDTO.ContentFile).
		Build()

	if err != nil {
		t.Errorf(UNEXPECTED_ERROR, err.Error())
		t.FailNow()
	}

	compareAllFieldsThumbnailImage(t, baseThumbnailDTO, image)
}

func compareAllFieldsThumbnailImage(t *testing.T, expected *thumbnailImageDTO.ThumbnailImageDTO, actual *thumbnailImageEntity.ThumbnailImage) {
	if expected.Id == nil && actual.GetId() != nil {
		t.Errorf("expected id nil, but got %v", actual.GetId())
	}
	if expected.Id != nil && actual.GetId() == nil {
		t.Errorf("expected id %v, but got nil", *expected.Id)
	} else if expected.Id != nil && *expected.Id != *actual.GetId() {
		t.Errorf("expected id %v, but got %v", *expected.Id, *actual.GetId())
	}

	compareCommonFieldsThumbnailImage(t, expected, actual)
}

func compareCommonFieldsThumbnailImage(t *testing.T, expected *thumbnailImageDTO.ThumbnailImageDTO, actual *thumbnailImageEntity.ThumbnailImage) {
	if expected.Name != actual.GetName() {
		t.Errorf("expected name %v, but got %v", expected.Name, actual.GetName())
	}
	if expected.Extension != actual.GetExtension() {
		t.Errorf("expected extension %v, but got %v", expected.Extension, actual.GetExtension())
	}
	if expected.Owner != actual.GetOwner() {
		t.Errorf("expected owner %v, but got %v", expected.Owner, actual.GetOwner())
	}
	if expected.Size != actual.GetSize() {
		t.Errorf("expected size %v, but got %v", expected.Size, actual.GetSize())
	}
	if expected.ContentFile != actual.GetContentFile() {
		t.Errorf("expected contentFile %v, but got %v", expected.ContentFile, actual.GetContentFile())
	}
}

func assertThumbnailBuilderException(t *testing.T, dto *thumbnailImageDTO.ThumbnailImageDTO, field string) {
	_, err := NewThumbnailImageBuilder().
		FromDTO(dto).
		Build()
	if err == nil {
		t.Errorf("Se esperaba un error al intentar crear una ThumbnailImage con el campo '%v' vacío", field)
		t.FailNow()
	}

	if err.Field != field {
		t.Errorf("Se esperaba un error del campo '%v', pero se obtuvo: %v", field, err.Field)
		t.FailNow()
	}
}

// Función para copiar el DTO base para cada caso de prueba
func copyThumbnailDTO() *thumbnailImageDTO.ThumbnailImageDTO {
	return &thumbnailImageDTO.ThumbnailImageDTO{
		Id:          baseThumbnailDTO.Id,
		Name:        baseThumbnailDTO.Name,
		Extension:   baseThumbnailDTO.Extension,
		ContentFile: baseThumbnailDTO.ContentFile,
		Owner:       baseThumbnailDTO.Owner,
		Size:        baseThumbnailDTO.Size,
	}
}
