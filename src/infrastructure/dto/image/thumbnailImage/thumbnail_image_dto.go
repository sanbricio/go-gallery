package thumbnailImageDTO

import thumbnailImageEntity "go-gallery/src/domain/entities/image/thumbnailImage"

// ThumbnailImageDTO representa la estructura de la imagen en miniatura
// @Description Contiene la información de la miniatura de una imagen, incluyendo su identificador, nombre, extensión, contenido en base64 y propietario (usuario)
type ThumbnailImageDTO struct {
	// Identificador de la miniatura
	Id *string `json:"id" bson:"_id,omitempty" example:"64a1f8b8e4b0c10d3c5b2e75"`

	// Identificador de la imagen
	ImageID *string `json:"imageID" bson:"imageID,omitempty" example:"64a1f8b8e4b0c20d3c5b2e90"`

	// Nombre del archivo de la miniatura
	Name string `json:"name" bson:"name" example:"prueba"`

	// Extensión del archivo de miniatura
	Extension string `json:"extension" bson:"extension" example:".jpeg"`

	// Contenido de la miniatura en base64
	ContentFile string `json:"content_file" bson:"content_file" example:"/9j/4AAQSkZJRgABAQEAAAAAAAD."`

	// Usuario propietario de la miniatura
	Owner string `json:"owner" bson:"owner" example:"usuario123"`

	// Tamaño de la miniatura en bytes
	Size string `json:"size" bson:"size" example:"2.3 kB"`

	// Tamaño de la imagen en bytes
	ImageSize string `json:"image_size" bson:"image_size" example:"2.3 kB"`
}

func FromThumbnailImage(thumbnailImage *thumbnailImageEntity.ThumbnailImage) *ThumbnailImageDTO {
	return &ThumbnailImageDTO{
		Id:          thumbnailImage.GetId(),
		Name:        thumbnailImage.GetName(),
		Extension:   thumbnailImage.GetExtension(),
		ContentFile: thumbnailImage.GetContentFile(),
		Owner:       thumbnailImage.GetOwner(),
		Size:        thumbnailImage.GetSize(),
		ImageSize:   thumbnailImage.GetImageSize(),
		ImageID:     thumbnailImage.GetImageID(),
	}
}
