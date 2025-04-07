package builder_test

import (
	"go-gallery/src/domain/entities/builder"
	"go-gallery/src/infrastructure/dto"
	"testing"
)

var baseDTO *dto.DTOImage

// Inicializamos el mock dto
func init() {
	baseDTO = &dto.DTOImage{
		IdImage:     "valid-id",
		Name:        "valid-name",
		Extension:   "jpg",
		ContentFile: "base64content",
		Owner:       "valid-owner",
		Size:        "1024",
	}
}

func TestImageBuilderEmptyFields(t *testing.T) {
	// Caso: 'name' vacío
	dto := copyDTO()
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

func assertBuilderException(t *testing.T, dto *dto.DTOImage, field string) {
	_, err := builder.NewImageBuilder().
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
func copyDTO() *dto.DTOImage {
	return &dto.DTOImage{
		IdImage:     baseDTO.IdImage,
		Name:        baseDTO.Name,
		Extension:   baseDTO.Extension,
		ContentFile: baseDTO.ContentFile,
		Owner:       baseDTO.Owner,
		Size:        baseDTO.Size,
	}
}
