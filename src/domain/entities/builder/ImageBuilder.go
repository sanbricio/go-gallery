package builder

import (
	"api-upload-photos/src/commons/exception"
	entity "api-upload-photos/src/domain/entities"
	"api-upload-photos/src/infrastructure/dto"

	"github.com/google/uuid"
)

type ImageBuilder struct {
	idImage     string
	name        string
	extension   string
	contentFile string
	owner       string
	size        string
}

func NewImageBuilder() *ImageBuilder {
	return &ImageBuilder{}
}

func (b *ImageBuilder) FromDTO(dto *dto.DTOImage) *ImageBuilder {
	b.idImage = dto.IdImage
	if b.idImage == "" {
		b.idImage = uuid.New().String()
	}
	b.name = dto.Name
	b.extension = dto.Extension
	b.contentFile = dto.ContentFile
	b.owner = dto.Owner
	b.size = dto.Size

	return b
}

func (b *ImageBuilder) Build() (*entity.Image, *exception.BuilderException) {
	if b.idImage == "" {
		return nil, exception.NewBuilderException("idImage", "El campo 'idImage' no debe estar vacio")
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

	return entity.NewImage(b.idImage, b.name, b.extension, b.contentFile, b.owner, b.size), nil
}

func (b *ImageBuilder) SetId(id string) *ImageBuilder {
	b.idImage = id
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
