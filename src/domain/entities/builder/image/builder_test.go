package imageBuilder

import (
	utilsImage "go-gallery/src/commons/utils/image"
	imageEntity "go-gallery/src/domain/entities/image"
	imageDTO "go-gallery/src/infrastructure/dto/image"
	"testing"
)

var baseDTO *imageDTO.ImageDTO

// Inicializamos el mock dto
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
	// Caso: 'id' vacío
	dto := copyDTO()
	dto.Id = nil
	assertBuilderException(t, dto, "id")

	// Caso: 'thumbnailId' vacío
	dto = copyDTO()
	dto.ThumbnailId = ""
	assertBuilderException(t, dto, "thumbnailId")

	// Caso: 'name' vacío
	dto = copyDTO()
	dto.Name = ""
	assertBuilderException(t, dto, "name")

	// Caso: 'extension' vacío
	dto = copyDTO()
	dto.Extension = ""
	assertBuilderException(t, dto, "extension")

	// Caso: 'contentFile' vacío
	dto = copyDTO()
	dto.ContentFile = ""
	assertBuilderException(t, dto, "contentFile")

	// Caso: 'owner' vacío
	dto = copyDTO()
	dto.Owner = ""
	assertBuilderException(t, dto, "owner")

	// Caso: 'size' vacío
	dto = copyDTO()
	dto.Size = ""
	assertBuilderException(t, dto, "size")
}

func TestImageBuilderNew(t *testing.T) {
	image, err := NewImageBuilder().
		SetThumbnailId(baseDTO.ThumbnailId).
		SetName(baseDTO.Name).
		SetExtension(baseDTO.Extension).
		SetOwner(baseDTO.Owner).
		SetSize(baseDTO.Size).
		SetContentFile(baseDTO.ContentFile).
		BuildNew()

	if err != nil {
		t.Errorf("No se esperaba un error al construir el la imagen: %v", err.Error())
		t.FailNow()
	}

	compareCommonFieldsImages(t, baseDTO, image)

	// Error a la hora de hacer un builder sin un campo especifico
	_, err = NewImageBuilder().
		SetName(baseDTO.Name).
		SetExtension(baseDTO.Extension).
		SetOwner(baseDTO.Owner).
		SetSize(baseDTO.Size).
		SetContentFile(baseDTO.ContentFile).
		BuildNew()

	if err == nil {
		t.Errorf("Se esperaba un error al construir el la imagen debido a que faltaba 'thumbnailID': %v", err.Error())
		t.FailNow()
	}

}

func TestImageBuilderFromImageUploadRequestDTO(t *testing.T) {
	rawContent := []byte("fake-binary-content")

	dto := &imageDTO.ImageUploadRequestDTO{
		Name:           "upload-name",
		Extension:      "png",
		RawContentFile: rawContent,
		Owner:          "upload-owner",
		Size:           "2048",
	}

	builder := NewImageBuilder().FromImageUploadRequestDTO(dto)
	image, err := builder.BuildNew()
	if err != nil {
		t.Errorf("No se esperaba un error al construir el la imagen: %v", err.Error())
		t.FailNow()
	}

	compareCommonFieldsImages(t, baseDTO, image)

	expectedContent := utilsImage.EncondeImageToBase64(rawContent)
	if image.GetContentFile() != expectedContent {
		t.Errorf("expected contentFile %v, but got %v", expectedContent, image.GetContentFile())
	}
}

func TestImageBuilderWithSetValues(t *testing.T) {
	image, err := NewImageBuilder().
		SetId(baseDTO.Id).
		SetThumbnailId(baseDTO.ThumbnailId).
		SetName(baseDTO.Name).
		SetExtension(baseDTO.Extension).
		SetOwner(baseDTO.Owner).
		SetSize(baseDTO.Size).
		SetContentFile(baseDTO.ContentFile).
		Build()

	if err != nil {
		t.Errorf("No se esperaba un error al construir el la imagen: %v", err.Error())
		t.FailNow()
	}

	compareAllFieldsImages(t, baseDTO, image)
}

func compareAllFieldsImages(t *testing.T, expected *imageDTO.ImageDTO, actual *imageEntity.Image) {
	if expected.Id != actual.GetId() {
		t.Errorf("expected id %v, but got %v", expected.Id, actual.GetId())
	}

	compareCommonFieldsImages(t, expected, actual)
}

func compareCommonFieldsImages(t *testing.T, expected *imageDTO.ImageDTO, actual *imageEntity.Image) {
	if expected.ThumbnailId != actual.GetThumbnailId() {
		t.Errorf("expected thumbnailId %v, but got %v", expected.ThumbnailId, actual.GetThumbnailId())
	}
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

func assertBuilderException(t *testing.T, dto *imageDTO.ImageDTO, field string) {
	_, err := NewImageBuilder().
		FromDTO(dto).
		Build()
	if err == nil {
		t.Errorf("Se esperaba un error al intentar crear una Image con el campo '%v' vacío", field)
		t.FailNow()
	}

	if err.Field != field {
		t.Errorf("Se esperaba un error del campo '%v', pero se obtuvo: %v", field, err)
		t.FailNow()
	}
}

// Función para copiar el DTO base para cada caso de prueba
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
