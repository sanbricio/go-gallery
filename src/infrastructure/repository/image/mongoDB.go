package image_repository

import (
	"context"
	"fmt"
	"go-gallery/src/commons/exception"

	imageDTO "go-gallery/src/infrastructure/dto/image"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const ImageMongoDBRepositoryKey = "ImageMongoDBRepository"

const (
	IMAGE_COLLECTION           string = "Image"
	THUMBNAIL_IMAGE_COLLECTION string = "ThumbnailImage "
	ID                         string = "id"
	OWNER                      string = "owner"
)

type ImageMongoDBRepository struct {
	mongoImage          *mongo.Collection
	mongoThumbnailImage *mongo.Collection
}

func NewImageMongoDBRepository(args map[string]string) ImageRepository {
	urlConnection := args["MONGODB_URL_CONNECTION"]
	databaseName := args["MONGODB_DATABASE"]

	db := connect(urlConnection, databaseName)

	repo := &ImageMongoDBRepository{
		mongoImage:          db.Collection(IMAGE_COLLECTION),
		mongoThumbnailImage: db.Collection(THUMBNAIL_IMAGE_COLLECTION),
	}

	return repo
}

func connect(urlConnection string, databaseName string) *mongo.Database {
	database, err := mongo.Connect(context.Background(), options.Client().ApplyURI(urlConnection))
	if err != nil {
		panic(fmt.Sprintf("No se ha podido conectar a MongoDB: %s", err.Error()))
	}

	err = database.Ping(context.Background(), readpref.Primary())
	if err != nil {
		panic(fmt.Sprintf("No se ha podido hacer ping a MongoDB: %s", err.Error()))
	}

	return database.Database(databaseName)
}

// TODO FIND all resized by owner

func (r *ImageMongoDBRepository) Find(dtoFind *imageDTO.ImageDTO) (*imageDTO.ImageDTO, *exception.ApiException) {
	filter := bson.M{
		ID:    dtoFind.Id,
		OWNER: dtoFind.Owner,
	}

	result, err := r.find(filter)
	if err != nil {
		return nil, err
	}
	return &result[0], nil
}

func (r *ImageMongoDBRepository) find(filter bson.M) ([]imageDTO.ImageDTO, *exception.ApiException) {
	cursor, err := r.mongoImage.Find(context.Background(), filter)
	if err != nil {
		return nil, exception.NewApiException(500, "Error al buscar las imagenes")
	}
	defer cursor.Close(context.Background())

	var results []imageDTO.ImageDTO
	for cursor.Next(context.Background()) {
		var image imageDTO.ImageDTO
		if err := cursor.Decode(&image); err != nil {
			return nil, exception.NewApiException(500, "Error al decodificar el las imagenes")
		}
		results = append(results, image)
	}

	if len(results) == 0 {
		return nil, exception.NewApiException(404, "Imagen no encontrada")
	}

	return results, nil
}

func (r *ImageMongoDBRepository) Insert(dtoInsertImage *imageDTO.ImageUploadRequestDTO) (*imageDTO.ImageUploadResponseDTO, *exception.ApiException) {
	// image, errBuilder := imageBuilder.NewImageBuilder().
	// 	FromDTO(dtoInsertImage).
	// 	Build()

	// if errBuilder != nil {
	// 	errorMessage := fmt.Sprintf("Error al construir la imagen: %s", errBuilder.Error())
	// 	return nil, exception.NewApiException(404, errorMessage)
	// }

	// _, err := r.mongoImage.InsertOne(context.Background(), image)
	// if err != nil {
	// 	return nil, exception.NewApiException(500, "Error al insertar el documento")
	// }

	// // resizedBytes, format, errResize := utils.ResizeImage(dtoInsertImage.RawContentFile, 200, 200)
	// // if errResize != nil {
	// // 	errorMessage := fmt.Sprintf("Error al generar miniatura: %s", errResize.Error())
	// // 	return nil, exception.NewApiException(500, errorMessage)
	// // }

	// return &imageDTO.ImageUploadResponseDTO{
	// 	Id:        image.GetId(),
	// 	Name:      image.GetName(),
	// 	Extension: image.GetExtension(),
	// 	Size:      image.GetSize(),
	// }, nil
	return nil, nil
}

func (r *ImageMongoDBRepository) Delete(dto *imageDTO.ImageDTO) (*imageDTO.ImageDTO, *exception.ApiException) {
	filter := bson.M{
		ID:    dto.Id,
		OWNER: dto.Owner,
	}
	foundImages, err := r.find(filter)
	if err != nil {
		return nil, err
	}

	_, errDelete := r.mongoImage.DeleteOne(context.Background(), filter)
	if errDelete != nil {
		return nil, exception.NewApiException(500, "Error al eliminar la imagen")
	}

	return &foundImages[0], nil
}
