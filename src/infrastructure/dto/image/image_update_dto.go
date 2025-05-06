package imageDTO

// ImageUpdateResponseDTO representa la peticion para actualizar la imagen.
type ImageUpdateRequestDTO struct {
	// ID de la imagen que queremos actualizar .
	Id string `json:"id" bson:"_id" example:"64a1f8b8e4b0c10d3c5b2e75"`

	// Nombre del archivo de imagen.
	Name string `json:"name" bson:"name" example:"foto_perfil"`

	// Usuario propietario de la imagen.
	Owner string `json:"owner" example:"usuario123"`

	// ID de la imagen miniatura asociada.s
	ThumbnailID string `json:"thumbnail_id" bson:"thumbnail_id" example:"64a1f8b8e4b0c10d3c5b2e75"`
}

// ImageUpdateResponseDTO representa la respuesta tras actualizar una imagen.
type ImageUpdateResponseDTO struct {
	// ID de la imagen actualizada.
	Id string `json:"id" example:"64a1f8b8e4b0c10d3c5b2e75"`

	// Usuario propietario de la imagen.
	Owner string `json:"owner" example:"usuario123"`

	// Campos que han sido actualizados en la imagen, representados como un mapa clave-valor.
	// Las claves son los nombres de los campos y los valores corresponden a los nuevos valores establecidos.
	UpdatedFields map[string]any `json:"updated_fields"`
}
