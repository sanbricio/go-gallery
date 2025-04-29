package thumbnailImageBuilder
// TODO REHACER ESTE TEST para thumbnail
// import (
// 	imageDTO "go-gallery/src/infrastructure/dto/image"
// 	"testing"
// )

// 
// var baseDTO *imageDTO.ImageDTO

// // Inicializamos el mock dto
// func init() {
// 	baseDTO = &imageDTO.ImageDTO{
// 		Id:          "valid-id",
// 		ThumbnailId: "valid-thumbnail-id",
// 		Name:        "valid-name",
// 		Extension:   "jpg",
// 		ContentFile: "base64content",
// 		Owner:       "valid-owner",
// 		Size:        "1024",
// 	}
// }

// func TestImageBuilderEmptyFields(t *testing.T) {
// 	// Caso: 'id' vacío
// 	dto := copyDTO()
// 	dto.Id = ""
// 	assertBuilderException(t, dto, "id")

// 	// Caso: 'thumbnailId' vacío
// 	dto = copyDTO()
// 	dto.ThumbnailId = ""
// 	assertBuilderException(t, dto, "thumbnailId")

// 	// Caso: 'name' vacío
// 	dto = copyDTO()
// 	dto.Name = ""
// 	assertBuilderException(t, dto, "name")

// 	// Caso: 'extension' vacío
// 	dto = copyDTO()
// 	dto.Extension = ""
// 	assertBuilderException(t, dto, "extension")

// 	// Caso: 'contentFile' vacío
// 	dto = copyDTO()
// 	dto.ContentFile = ""
// 	assertBuilderException(t, dto, "contentFile")

// 	// Caso: 'owner' vacío
// 	dto = copyDTO()
// 	dto.Owner = ""
// 	assertBuilderException(t, dto, "owner")

// 	// Caso: 'size' vacío
// 	dto = copyDTO()
// 	dto.Size = ""
// 	assertBuilderException(t, dto, "size")
// }

// func assertBuilderException(t *testing.T, dto *imageDTO.ImageDTO, field string) {
// 	_, err := NewImageBuilder().
// 		FromDTO(dto).
// 		Build()
// 	if err == nil {
// 		t.Errorf("Se esperaba un error al intentar crear una Image con el campo '%v' vacío", field)
// 		t.FailNow()
// 	}

// 	if err.Field != field {
// 		t.Errorf("Se esperaba un error del campo '%v', pero se obtuvo: %v", field, err)
// 		t.FailNow()
// 	}
// }

// // Función para copiar el DTO base para cada caso de prueba
// func copyDTO() *imageDTO.ImageDTO {
// 	return &imageDTO.ImageDTO{
// 		Id:          baseDTO.Id,
// 		ThumbnailId: baseDTO.ThumbnailId,
// 		Name:        baseDTO.Name,
// 		Extension:   baseDTO.Extension,
// 		ContentFile: baseDTO.ContentFile,
// 		Owner:       baseDTO.Owner,
// 		Size:        baseDTO.Size,
// 	}
// }
