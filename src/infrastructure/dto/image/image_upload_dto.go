package imageDTO

// ImageUploadRequestDTO representa la estructura de la petición para subir una imagen.
// Contiene los metadatos de la imagen y el contenido en diferentes formatos.
type ImageUploadRequestDTO struct {
	// Nombre del archivo de imagen.
	// Example: foto_perfil
	Name string `json:"name" bson:"name" example:"foto_perfil"`

	// Extensión del archivo de imagen (ejemplo: .jpeg, .png).
	// Example: .jpeg
	Extension string `json:"extension" bson:"extension" example:".jpeg"`

	// Contenido de la imagen codificado en base64.
	// Nota: utilizado si se envía en texto plano.
	// Example: /9j/4AAQSkZJRgABAQEAAAAAAAD...
	ContentFile string `json:"content_file" bson:"content_file" example:"/9j/4AAQSkZJRgABAQEAAAAAAAD..."`

	// Contenido de la imagen en bytes sin codificar.
	// Example: []byte
	RawContentFile []byte `json:"raw_content_file" bson:"raw_content_file"`

	// Usuario propietario de la imagen.
	// Example: usuario123
	Owner string `json:"owner" bson:"owner" example:"usuario123"`

	// Tamaño del archivo de imagen como string.
	// Example: 204800
	Size string `json:"size" bson:"size" example:"204800"`
}

// ImageUploadResponseDTO representa la respuesta tras subir una imagen.
// Devuelve los principales metadatos de la imagen almacenada.
type ImageUploadResponseDTO struct {
	// ID de la imagen almacenada.
	// Example: 64a1f8b8e4b0c10d3c5b2e75
	Id string `json:"id" bson:"id" example:"64a1f8b8e4b0c10d3c5b2e75"`

	// ID de la imagen miniatura asociada.
	// Example: 64a1f8b8e4b0c10d3c5b2e75
	ThumbnailId string `json:"thumbnail_id" bson:"thumbnail_id" example:"64a1f8b8e4b0c10d3c5b2e75"`

	// Nombre de la imagen almacenada.
	// Example: foto_perfil
	Name string `json:"name" bson:"name" example:"foto_perfil"`

	// Extensión de la imagen almacenada.
	// Example: .jpeg
	Extension string `json:"extension" bson:"extension" example:".jpeg"`

	// Tamaño de la imagen almacenada en bytes como string.
	// Example: 204800
	Size string `json:"size" bson:"size" example:"204800"`
}
