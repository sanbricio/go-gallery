package repository

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

const (
	IMAGE_COLLECTION = "Images"
	ID_IMAGE         = "id_image"
)

type RepositoryImageMongoDB struct {
	client *mongo.Database
}

func NewRepositoryImageMongoDB(urlConnection string, databaseName string) (*RepositoryImageMongoDB, *exception.ConnectionException) {
	db, err := connect(urlConnection, databaseName)
	if err != nil {
		return nil, err
	}

	repo := &RepositoryImageMongoDB{
		client: db,
	}
	return repo, nil
}

func connect(urlConnection string, databaseName string) (*mongo.Database, *exception.ConnectionException) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(urlConnection))
	if err != nil {
		return nil, exception.NewConnectionException(fmt.Sprintf("No se ha podido conectar a MongoDB: %s", err.Error()), err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, exception.NewConnectionException(fmt.Sprintf("No se ha podido hacer ping a MongoDB: %s", err.Error()), err)
	}

	database := client.Database(databaseName)
	return database, nil
}

func (r *RepositoryImageMongoDB) Find(id string) (*dto.DTOImage, *exception.ApiException) {
	collection := r.client.Collection(IMAGE_COLLECTION)
	var result dto.DTOImage
	filter := bson.M{ID_IMAGE: id}

	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, exception.NewApiException(404, "No se encontró el documento en base de datos")
		}
		return nil, exception.NewApiException(500, "Error al buscar el documento")
	}

	return &result, nil
}

func (r *RepositoryImageMongoDB) Insert(processedImage *dto.DTOProcessedImage) (*dto.DTOImage, *exception.ApiException) {
	collection := r.client.Collection(IMAGE_COLLECTION)

	image, errBuilder := builder.NewImageBuilder().
		From(processedImage.FileName, processedImage.FileExtension, processedImage.EncodedData, "SANTI", processedImage.FileSizeHumanReadable).
		Build()

	if errBuilder != nil {
		errorMessage := fmt.Sprintf("Error al construir la imagen: %s", errBuilder.Error())
		return nil, exception.NewApiException(404, errorMessage)
	}

	dto := dto.FromImage(image)
	_, err := collection.InsertOne(context.Background(), dto)
	if err != nil {
		return nil, exception.NewApiException(500, "Error al insertar el documento")
	}

	return dto, nil
}

func (r *RepositoryImageMongoDB) Delete(id string) (*dto.DTOImage, *exception.ApiException) {
	collection := r.client.Collection(IMAGE_COLLECTION)
	var result dto.DTOImage
	filter := bson.M{ID_IMAGE: id}

	err := collection.FindOneAndDelete(context.Background(), filter)
	if err != nil {
		if err.Err() == mongo.ErrNoDocuments {
			return nil, exception.NewApiException(404, "No se encontró el documento en base de datos")
		}

		return nil, exception.NewApiException(500, "Error al eliminar el documento")
	}

	return &result, nil
}
