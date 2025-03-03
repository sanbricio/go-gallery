package image_repository

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/domain/entities/builder"
	"api-upload-photos/src/infrastructure/dto"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const ImageMongoDBRepositoryKey = "ImageMongoDBRepository"

const (
	IMAGE_COLLECTION = "Images"
	ID_IMAGE         = "id_image"
	OWNER            = "owner"
)

type ImageMongoDBRepository struct {
	mongo *mongo.Collection
}

func NewImageMongoDBRepository(args map[string]string) ImageRepository {
	urlConnection := args["MONGODB_URL_CONNECTION"]
	databaseName := args["MONGODB_DATABASE"]

	db := connect(urlConnection, databaseName)

	repo := &ImageMongoDBRepository{
		mongo: db.Collection(IMAGE_COLLECTION),
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

func (r *ImageMongoDBRepository) Find(dtoFind *dto.DTOImage) (*dto.DTOImage, *exception.ApiException) {
	filter := bson.M{
		ID_IMAGE: dtoFind.IdImage,
		OWNER:    dtoFind.Owner,
	}

	result, err := r.find(filter)
	if err != nil {
		return nil, err
	}
	return &result[0], nil
}

func (r *ImageMongoDBRepository) find(filter bson.M) ([]dto.DTOImage, *exception.ApiException) {
	cursor, err := r.mongo.Find(context.Background(), filter)
	if err != nil {
		return nil, exception.NewApiException(500, "Error al buscar las imagenes")
	}
	defer cursor.Close(context.Background())

	var results []dto.DTOImage
	for cursor.Next(context.Background()) {
		var image dto.DTOImage
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

func (r *ImageMongoDBRepository) Insert(dtoInsertImage *dto.DTOImage) (*dto.DTOImage, *exception.ApiException) {
	image, errBuilder := builder.NewImageBuilder().
		FromDTO(dtoInsertImage).
		Build()

	if errBuilder != nil {
		errorMessage := fmt.Sprintf("Error al construir la imagen: %s", errBuilder.Error())
		return nil, exception.NewApiException(404, errorMessage)
	}

	dto := dto.FromImage(image)
	_, err := r.mongo.InsertOne(context.Background(), dto)
	if err != nil {
		return nil, exception.NewApiException(500, "Error al insertar el documento")
	}

	return dto, nil
}

func (r *ImageMongoDBRepository) Delete(dto *dto.DTOImage) (*dto.DTOImage, *exception.ApiException) {
	filter := bson.M{
		ID_IMAGE: dto.IdImage,
		OWNER:    dto.Owner,
	}
	foundImages, err := r.find(filter)
	if err != nil {
		return nil, err
	}

	_, errDelete := r.mongo.DeleteOne(context.Background(), filter)
	if errDelete != nil {
		return nil, exception.NewApiException(500, "Error al eliminar la imagen")
	}

	return &foundImages[0], nil
}
