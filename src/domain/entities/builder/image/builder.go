package imageBuilder

import (
	"go-gallery/src/commons/exception"
	imageEntity "go-gallery/src/domain/entities/image"
	imageDTO "go-gallery/src/infrastructure/dto/image"
)

type ImageBuilder struct {
	id          string
	thumbnailId string
	name        string
	extension   string
	contentFile string
	owner       string
	size        string
}

func NewImageBuilder() *ImageBuilder {
	return &ImageBuilder{}
}

func (b *ImageBuilder) FromDTO(dto *imageDTO.ImageDTO) *ImageBuilder {
	b.id = dto.Id
	b.thumbnailId = dto.ThumbnailId
	b.name = dto.Name
	b.extension = dto.Extension
	b.contentFile = dto.ContentFile
	b.owner = dto.Owner
	b.size = dto.Size

	return b
}

func (b *ImageBuilder) Build() (*imageEntity.Image, *exception.BuilderException) {
	if b.id == "" {
		return nil, exception.NewBuilderException("id", "El campo 'id' no debe estar vacio")
	}

	if b.thumbnailId == "" {
		return nil, exception.NewBuilderException("thumbnailId", "El campo 'thumbnailId' no debe estar vacio")
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

	return imageEntity.NewImage(b.id, b.thumbnailId, b.name, b.extension, b.contentFile, b.owner, b.size), nil
}

func (b *ImageBuilder) SetId(id string) *ImageBuilder {
	b.id = id
	return b
}

func (b *ImageBuilder) SetThumbnailId(thumbnailId string) *ImageBuilder {
	b.thumbnailId = thumbnailId
	return b
}

func (b *ImageBuilder) SetName(name string) *ImageBuilder {
	b.name = name
	return b
}

func (b *ImageBuilder) SetExtension(extension string) *ImageBuilder {
	b.extension = extension
	return b
}

func (b *ImageBuilder) SetContentFile(contentFile string) *ImageBuilder {
	b.contentFile = contentFile
	return b
}

func (b *ImageBuilder) SetOwner(owner string) *ImageBuilder {
	b.owner = owner
	return b
}

func (b *ImageBuilder) SetSize(size string) *ImageBuilder {
	b.size = size
	return b
}
