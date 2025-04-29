package imageDTO

import (
	imageEntity "go-gallery/src/domain/entities/image"
)

// ImageDTO representa la estructura de una imagen
// @Description Contiene la información de una imagen, incluyendo su identificador, nombre, extensión, contenido en base64 y propietario (usuario)
type ImageDTO struct {
	// ID de la imagen
	// Example: 64a1f8b8e4b0c10d3c5b2e75
	Id *string `json:"id" bson:"_id,omitempty" example:"64a1f8b8e4b0c10d3c5b2e75"`

	// ID de la imagen miniatura asociada.
	// Example: 64a1f8b8e4b0c10d3c5b2e75
	ThumbnailId string `json:"thumbnail_id" bson:"thumbnail_id" example:"64a1f8b8e4b0c10d3c5b2e75"`

	// Nombre del archivo de la imagen
	// Example: foto_perfil
	Name string `json:"name" bson:"name" example:"prueba"`

	// Extensión del archivo de imagen
	// Example: jpg
	Extension string `json:"extension" bson:"extension" example:".jpeg"`

	// Contenido de la imagen en base64
	// Example: /9j/4AAQSkZJRgABAQEAAAAAAAD...
	ContentFile string `json:"content_file" bson:"content_file" example:"/9j/4AAQSkZJRgABAQEAAAAAAAD."`

	// Usuario propietario de la imagen
	// Example: usuario123
	Owner string `json:"owner" bson:"owner" example:"usuario123"`

	// Tamaño de la imagen en bytes
	// Example: 204800
	Size string `json:"size" bson:"size" example:"2.3 kB"`
}

func FromImage(image *imageEntity.Image) *ImageDTO {
	id := image.GetId()
	return &ImageDTO{
		Id:          id,
		ThumbnailId: image.GetThumbnailId(),
		Name:        image.GetName(),
		Extension:   image.GetExtension(),
		ContentFile: image.GetContentFile(),
		Owner:       image.GetOwner(),
		Size:        image.GetSize(),
	}
}
