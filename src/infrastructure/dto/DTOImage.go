package dto

import entity "api-upload-photos/src/domain/entities"

// DTOImage representa la estructura de una imagen
// @Description Contiene la información de una imagen, incluyendo su identificador, nombre, extensión, contenido en base64 y propietario (usuario)
type DTOImage struct {
	// ID único de la imagen
	// Example: 64a1f8b8e4b0c10d3c5b2e75
	IdImage string `json:"id_image" bson:"id_image" example:"64a1f8b8e4b0c10d3c5b2e75"`

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

func FromImage(image *entity.Image) *DTOImage {
	return &DTOImage{
		IdImage:     image.GetId(),
		Name:        image.GetName(),
		Extension:   image.GetExtension(),
		ContentFile: image.GetContentFile(),
		Owner:       image.GetOwner(),
		Size:        image.GetSize(),
	}
}
