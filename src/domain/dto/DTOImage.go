package dto

import entity "api-upload-photos/src/domain/entities"

type DTOImage struct {
	IdImage     string `json:"id_image" bson:"id_image"`
	Name        string `json:"name" bson:"name"`
	Extension   string `json:"extension" bson:"extension"`
	ContentFile string `json:"content_file" bson:"content_file"`
	Owner       string `json:"owner" bson:"owner"`
	Size        string `json:"size" bson:"size"`
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

func (img *DTOImage) AsImageEntity() *entity.Image {
	return entity.NewImageFromDTO(img.IdImage, img.Name, img.Extension, img.ContentFile, img.Owner, img.Size)
}
