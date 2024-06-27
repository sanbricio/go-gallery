package infrastructure

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/domain/dto"
	entity "api-upload-photos/src/domain/entities"
	handler "api-upload-photos/src/infrastructure"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	ImagesCollection = "Images"
	KIdImage         = "id_image"
)

type RepositoryMongoDB struct {
	client *mongo.Database
}

func NewRepositoryMongoDB(urlConnection string, databaseName string) (*RepositoryMongoDB, *exception.ApiException) {
	db, err := connect(urlConnection, databaseName)
	if err != nil {
		return nil, err
	}

	repo := &RepositoryMongoDB{
		client: db,
	}
	return repo, nil
}

// TODO Tratar un nuevo tipo de excepcion ej: ConnectionException
func connect(urlConnection string, databaseName string) (*mongo.Database, *exception.ApiException) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(urlConnection))
	if err != nil {
		return nil, exception.NewApiException(500, fmt.Sprintf("No se ha podido conectar a MongoDB: %s", err.Error()))
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, exception.NewApiException(500, fmt.Sprintf("No se ha podido hacer ping a MongoDB: %s", err.Error()))
	}

	database := client.Database(databaseName)
	return database, nil
}

func (r *RepositoryMongoDB) Find(id string) (*dto.DTOImage, *exception.ApiException) {
	collection := r.client.Collection(ImagesCollection)
	var result dto.DTOImage
	filter := bson.M{KIdImage: id}

	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, exception.NewApiException(404, "No se encontró el documento en base de datos")
		}
		return nil, exception.NewApiException(500, "Error al buscar el documento")
	}

	return &result, nil
}

// TODO Probar cuando este terminado desarrollo front
func (r *RepositoryMongoDB) Insert(processedImage *handler.ProcessedImage) (*dto.DTOImage, *exception.ApiException) {
	collection := r.client.Collection(ImagesCollection)

	image := entity.NewImage(processedImage.FileName, processedImage.FileExtension, processedImage.EncodedData, "SANTI", processedImage.FileSizeHumanReadable)

	dto := dto.FromImage(image)
	_, err := collection.InsertOne(context.Background(), dto)
	if err != nil {
		return nil, exception.NewApiException(500, "Error al insertal el documento")
	}

	return dto, nil
}

func (r *RepositoryMongoDB) Delete(id string) (*dto.DTOImage, *exception.ApiException) {
	collection := r.client.Collection(ImagesCollection)
	var result dto.DTOImage
	filter := bson.M{KIdImage: id}

	err := collection.FindOneAndDelete(context.Background(), filter)
	if err != nil {
		if err.Err() == mongo.ErrNoDocuments {
			return nil, exception.NewApiException(404, "No se encontró el documento en base de datos")
		}

		return nil, exception.NewApiException(500, "Error al eliminar el documento")
	}

	return &result, nil
}
