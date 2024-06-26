package dto

import entity "api-upload-photos/src/domain/entities"

type DTOImage struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Extension   string `json:"extension"`
	ContentFile string `json:"content_file"`
	Owner       string `json:"owner"`
	Size        string `json:"size"`
}

func FromImage(image *entity.Image) *DTOImage {
	return &DTOImage{
		Id:          image.GetId(),
		Name:        image.GetName(),
		Extension:   image.GetExtension(),
		ContentFile: image.GetContentFile(),
		Owner:       image.GetOwner(),
		Size:        image.GetSize(),
	}
}

func (img *DTOImage) AsImageEntity() *entity.Image {
	return entity.NewImageFromDTO(img.Id, img.Name, img.Extension, img.ContentFile, img.Owner, img.Size)
}
