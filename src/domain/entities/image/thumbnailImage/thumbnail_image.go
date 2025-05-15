package thumbnailImageEntity

type ThumbnailImage struct {
	id          *string
	imageID     *string
	name        string
	extension   string
	contentFile string
	size        string
	owner       string
	imageSize   string
}

func NewThumbnailImage(id, imageID *string, name, extension, contentFile, size, owner, imageSize string) *ThumbnailImage { // NOSONAR
	return &ThumbnailImage{
		id:          id,
		imageID:     imageID,
		name:        name,
		extension:   extension,
		contentFile: contentFile,
		size:        size,
		owner:       owner,
		imageSize:   imageSize,
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

func (img *ThumbnailImage) GetImageID() *string {
	return img.imageID
}

func (img *ThumbnailImage) GetImageSize() string {
	return img.imageSize
}
