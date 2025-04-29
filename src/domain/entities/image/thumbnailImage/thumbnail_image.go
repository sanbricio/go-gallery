package thumbnailImageEntity

type ThumbnailImage struct {
	id          *string
	name        string
	extension   string
	contentFile string
	size        string
	owner       string
}

func NewThumbnailImage(id *string, name, extension, contentFile, size, owner string) *ThumbnailImage {
	return &ThumbnailImage{
		id:          id,
		name:        name,
		extension:   extension,
		contentFile: contentFile,
		size:        size,
		owner:       owner,
	}
}

func (img *ThumbnailImage) GetId() *string {
	return img.id
}

func (img *ThumbnailImage) GetName() string {
	return img.name
}

func (img *ThumbnailImage) GetExtension() string {
	return img.extension
}

func (img *ThumbnailImage) GetContentFile() string {
	return img.contentFile
}

func (img *ThumbnailImage) GetSize() string {
	return img.size
}

func (img *ThumbnailImage) GetOwner() string {
	return img.owner
}
