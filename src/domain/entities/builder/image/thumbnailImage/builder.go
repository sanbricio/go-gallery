package thumbnailImageBuilder

import (
	"go-gallery/src/commons/exception"
	utilsImage "go-gallery/src/commons/utils/image"
	thumbnailImageEntity "go-gallery/src/domain/entities/image/thumbnailImage"
	imageDTO "go-gallery/src/infrastructure/dto/image"
	thumbnailImageDTO "go-gallery/src/infrastructure/dto/image/thumbnailImage"
)

type ThumbnailImageBuilder struct {
	id          *string
	name        string
	extension   string
	contentFile string
	owner       string
	size        string
}

func NewThumbnailImageBuilder() *ThumbnailImageBuilder {
	return &ThumbnailImageBuilder{}
}

func (b *ThumbnailImageBuilder) FromImageUploadRequestDTO(dto *imageDTO.ImageUploadRequestDTO) *ThumbnailImageBuilder {
	b.name = dto.Name
	b.extension = dto.Extension
	b.owner = dto.Owner
	b.size = dto.Size

	return b
}

func (b *ThumbnailImageBuilder) FromDTO(dto *thumbnailImageDTO.ThumbnailImageDTO) *ThumbnailImageBuilder {
	id := dto.Id
	b.id = id
	b.name = dto.Name
	b.extension = dto.Extension
	b.contentFile = dto.ContentFile
	b.owner = dto.Owner
	b.size = dto.Size

	return b
}

func (b *ThumbnailImageBuilder) SetId(id *string) *ThumbnailImageBuilder {
	b.id = id
	return b
}

func (b *ThumbnailImageBuilder) SetName(name string) *ThumbnailImageBuilder {
	b.name = name
	return b
}

func (b *ThumbnailImageBuilder) SetExtension(extension string) *ThumbnailImageBuilder {
	b.extension = extension
	return b
}

func (b *ThumbnailImageBuilder) SetContentFile(contentFile []byte) *ThumbnailImageBuilder {
	b.contentFile = utilsImage.EncondeImageToBase64(contentFile)
	return b
}

func (b *ThumbnailImageBuilder) SetSize(size string) *ThumbnailImageBuilder {
	b.size = size
	return b
}

func (b *ThumbnailImageBuilder) BuildNew() (*thumbnailImageEntity.ThumbnailImage, *exception.BuilderException) {
	if b.name == "" {
		return nil, exception.NewBuilderException("name", "El campo 'name' no debe estar vacío")
	}
	if b.extension == "" {
		return nil, exception.NewBuilderException("extension", "El campo 'extension' no debe estar vacío")
	}
	if b.contentFile == "" {
		return nil, exception.NewBuilderException("contentFile", "El campo 'contentFile' no debe estar vacío")
	}
	if b.size == "" {
		return nil, exception.NewBuilderException("size", "El campo 'size' no debe estar vacío")
	}
	if b.owner == "" {
		return nil, exception.NewBuilderException("owner", "El campo 'owner' no debe estar vacio")
	}

	return thumbnailImageEntity.NewThumbnailImage(nil, b.name, b.extension, b.contentFile, b.size, b.owner), nil
}

func (b *ThumbnailImageBuilder) Build() (*thumbnailImageEntity.ThumbnailImage, *exception.BuilderException) {
	if b.id == nil {
		return nil, exception.NewBuilderException("id", "El campo 'id' no debe estar vacio")
	}

	if b.name == "" {
		return nil, exception.NewBuilderException("name", "El campo 'name' no debe estar vacio")
	}

	if b.extension == "" {
		return nil, exception.NewBuilderException("extension", "El campo 'extension' no debe de estar vacio")
	}

	if b.contentFile == "" {
		return nil, exception.NewBuilderException("contentFile", "El campo 'contentFile' no debe de estar vacio")
	}

	if b.owner == "" {
		return nil, exception.NewBuilderException("owner", "El campo 'owner' no debe estar vacio")
	}

	if b.size == "" {
		return nil, exception.NewBuilderException("size", "El campo 'size' no debe estar vacio")
	}

	return thumbnailImageEntity.NewThumbnailImage(b.id, b.name, b.extension, b.contentFile, b.size, b.owner), nil
}
