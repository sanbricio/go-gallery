package entity

type Image struct {
	idImage     string
	name        string
	extension   string
	contentFile string
	owner       string
	size        string
}

func NewImage(idImage, name, extension, contentFile, owner, size string) *Image {
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
