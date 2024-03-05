package domain

type DTOImage struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Extension   string `json:"extension"`
	ContentFile string `json:"content_file"`
	Owner       string `json:"owner"`
	Size        string `json:"size"`
}

func NewDTOImage(id string, name string, extension string, contentFile string, owner string, size string) *DTOImage {
	return &DTOImage{
		Id:          id,
		Name:        name,
		Extension:   extension,
		ContentFile: contentFile,
		Owner:       owner,
		Size:        size,
	}
}

func FromDTO(dto *DTOImage) *Image {
	return NewImage(dto.Id, dto.Name, dto.Extension, dto.ContentFile, dto.Owner, dto.Size)
}
