package imageBuilder

import (
	"go-gallery/src/commons/exception"
	utilsImage "go-gallery/src/commons/utils/image"
	validators "go-gallery/src/commons/utils/validations"
	imageEntity "go-gallery/src/domain/entities/image"
	imageDTO "go-gallery/src/infrastructure/dto/image"
)

type ImageBuilder struct {
	id          *string
	name        string
	extension   string
	contentFile string
	owner       string
	size        string
}

func NewImageBuilder() *ImageBuilder {
	return &ImageBuilder{}
}

func (b *ImageBuilder) FromImageUploadRequestDTO(dto *imageDTO.ImageUploadRequestDTO) *ImageBuilder {
	b.name = dto.Name
	b.extension = dto.Extension
	b.contentFile = utilsImage.EncondeImageToBase64(dto.RawContentFile)
	b.owner = dto.Owner
	b.size = dto.Size

	return b
}

func (b *ImageBuilder) FromDTO(dto *imageDTO.ImageDTO) *ImageBuilder {
	id := dto.Id
	b.id = id
	b.name = dto.Name
	b.extension = dto.Extension
	b.contentFile = dto.ContentFile
	b.owner = dto.Owner
	b.size = dto.Size

	return b
}

func (b *ImageBuilder) BuildNew() (*imageEntity.Image, *exception.BuilderException) {
	err := b.validateCommons()
	if err != nil {
		return nil, err
	}

	return imageEntity.NewImage(nil, b.name, b.extension, b.contentFile, b.owner, b.size), nil
}

func (b *ImageBuilder) Build() (*imageEntity.Image, *exception.BuilderException) {
	err := b.validateAll()
	if err != nil {
		return nil, err
	}

	return imageEntity.NewImage(b.id, b.name, b.extension, b.contentFile, b.owner, b.size), nil
}

func (b *ImageBuilder) validateAll() *exception.BuilderException {
	err := validators.ValidateNonEmptyStringField("id", b.id)
	if err != nil {
		return exception.NewBuilderException("id", err.Error())
	}

	errBuilder := b.validateCommons()
	if errBuilder != nil {
		return errBuilder
	}

	return nil
}

func (b *ImageBuilder) validateCommons() *exception.BuilderException {

	if err := validators.ValidateNonEmptyStringField("name", b.name); err != nil {
		return exception.NewBuilderException("name", err.Error())
	}

	if err := validators.ValidateNonEmptyStringField("extension", b.extension); err != nil {
		return exception.NewBuilderException("extension", err.Error())
	}

	if err := validators.ValidateNonEmptyStringField("contentFile", b.contentFile); err != nil {
		return exception.NewBuilderException("contentFile", err.Error())
	}

	if err := validators.ValidateNonEmptyStringField("owner", b.owner); err != nil {
		return exception.NewBuilderException("owner", err.Error())
	}

	if err := validators.ValidateNonEmptyStringField("size", b.size); err != nil {
		return exception.NewBuilderException("size", err.Error())
	}
	return nil
}

func (b *ImageBuilder) SetId(id *string) *ImageBuilder {
	b.id = id
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
