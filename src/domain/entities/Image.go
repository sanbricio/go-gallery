package entity

import "github.com/google/uuid"

type Image struct {
	idImage     string
	name        string
	extension   string
	contentFile string
	owner       string
	size        string
}

func NewImage(name string, extension string, contentFile string, owner string, size string) *Image {
	return &Image{
		idImage:     uuid.New().String(),
		name:        name,
		extension:   extension,
		contentFile: contentFile,
		owner:       owner,
		size:        size,
	}
}

func NewImageFromDTO(idImage string, name string, extension string, contentFile string, owner string, size string) *Image {
	return &Image{
		idImage:     idImage,
		name:        name,
		extension:   extension,
		contentFile: contentFile,
		owner:       owner,
		size:        size,
	}
}

func (img *Image) GetId() string {
	return img.idImage
}

func (img *Image) GetName() string {
	return img.name
}

func (img *Image) GetExtension() string {
	return img.extension
}

func (img *Image) GetContentFile() string {
	return img.contentFile
}

func (img *Image) GetOwner() string {
	return img.owner
}

func (img *Image) GetSize() string {
	return img.size
}

func (img *Image) SetName(name string) {
	img.name = name
}

func (img *Image) SetExtension(extension string) {
	img.extension = extension
}

func (img *Image) SetContentFile(contentFile string) {
	img.contentFile = contentFile
}

func (img *Image) SetOwner(owner string) {
	img.owner = owner
}
