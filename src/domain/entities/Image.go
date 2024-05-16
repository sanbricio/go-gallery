package entity

type Image struct {
	id          string
	name        string
	extension   string
	contentFile string
	owner       string
	size        string
}

func NewImage(id string, name string, extension string, contentFile string, owner string, size string) *Image {
	return &Image{
		id:          id,
		name:        name,
		extension:   extension,
		contentFile: contentFile,
		owner:       owner,
		size:        size,
	}
}

func (img *Image) GetId() string {
	return img.id
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