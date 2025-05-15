package thumbnailImageDTO

// ThumbnailImageCursorDTO representa un cursor de paginación para las miniaturas
// @Description Contiene una lista de miniaturas y el identificador del último elemento para la paginación.
type ThumbnailImageCursorDTO struct {
	// Lista de miniaturas de imagen
	Thumbnails []ThumbnailImageDTO `json:"thumbnails" bson:"thumbnails"`
	// ID del último elemento para la paginación
	LastID string `json:"lastID" bson:"lastID,omitempty" example:"64a1f8b8e4b0c10d3c5b2e75"`
}
