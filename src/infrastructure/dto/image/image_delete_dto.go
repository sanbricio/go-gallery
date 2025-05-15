package imageDTO

// ImageUpdateResponseDTO representa la peticion para eliminar una imagen.
type ImageDeleteRequestDTO struct {
	// ID de la imagen que queremos actualizar .
	Id string `json:"id" bson:"_id" example:"64a1f8b8e4b0c10d3c5b2e75"`

	// Usuario propietario de la imagen.
	Owner string `json:"owner" example:"usuario123"`

	// ID de la imagen miniatura asociadas
	ThumbnailID string `json:"thumbnail_id" bson:"thumbnail_id" example:"64a1f8b8e4b0c10d3c5b2e75"`
}
