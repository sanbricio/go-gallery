package imageBuilder

import (
	"go-gallery/src/commons/exception"
	utilsImage "go-gallery/src/commons/utils/image"
	imageEntity "go-gallery/src/domain/entities/image"
	imageDTO "go-gallery/src/infrastructure/dto/image"
)

type ImageBuilder struct {
	id          *string
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
	b.thumbnailId = dto.ThumbnailId
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

	return imageEntity.NewImage(nil, b.thumbnailId, b.name, b.extension, b.contentFile, b.owner, b.size), nil
}

func (b *ImageBuilder) Build() (*imageEntity.Image, *exception.BuilderException) {
	err := b.validateAll()
	if err != nil {
		return nil, err
	}

	return imageEntity.NewImage(b.id, b.thumbnailId, b.name, b.extension, b.contentFile, b.owner, b.size), nil
}

func (b *ImageBuilder) validateAll() *exception.BuilderException {
	err := b.validateStringPointerField("id", b.id)
	if err != nil {
		return err
	}

	err = b.validateCommons()
	if err != nil {
		return err
	}

	return nil
}

func (b *ImageBuilder) validateCommons() *exception.BuilderException {
	if err := b.validateStringField("thumbnailId", b.thumbnailId); err != nil {
		return err
	}

	if err := b.validateStringField("name", b.name); err != nil {
		return err
	}

	if err := b.validateStringField("extension", b.extension); err != nil {
		return err
	}
	if err := b.validateStringField("contentFile", b.contentFile); err != nil {
		return err
	}
	if err := b.validateStringField("owner", b.owner); err != nil {
		return err
	}
	if err := b.validateStringField("size", b.size); err != nil {
		return err
	}
	return nil
}

func (b *ImageBuilder) validateStringField(fieldName, fieldValue string) *exception.BuilderException {
	if fieldValue == "" {
		return exception.NewBuilderException(fieldName, "El campo '"+fieldName+"' no debe estar vacio")
	}
	return nil
}

func (b *ImageBuilder) validateStringPointerField(fieldName string, fieldValue *string) *exception.BuilderException {
	if fieldValue == nil || *fieldValue == "" {
		return exception.NewBuilderException(fieldName, "El campo '"+fieldName+"' no debe estar vacio")
	}
	return nil
}

func (b *ImageBuilder) SetId(id *string) *ImageBuilder {
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
