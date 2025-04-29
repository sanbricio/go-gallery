package imageRepository

import (
	"context"
	"fmt"
	"go-gallery/src/commons/exception"

	imageBuilder "go-gallery/src/domain/entities/builder/image"

	imageDTO "go-gallery/src/infrastructure/dto/image"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const ImageMongoDBRepositoryKey = "ImageMongoDBRepository"

const (
	IMAGE_COLLECTION string = "Image"
	ID               string = "_id"
	OWNER            string = "owner"
)

type ImageMongoDBRepository struct {
	mongoImage *mongo.Collection
}

func NewImageMongoDBRepository(args map[string]string) ImageRepository {
	urlConnection := args["MONGODB_URL_CONNECTION"]
	databaseName := args["MONGODB_DATABASE"]

	db := connect(urlConnection, databaseName)

	repo := &ImageMongoDBRepository{
		mongoImage: db.Collection(IMAGE_COLLECTION),
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

func (r *ImageMongoDBRepository) Insert(dtoInsertImage *imageDTO.ImageUploadRequestDTO, thumbnailImageID string) (*imageDTO.ImageUploadResponseDTO, *exception.ApiException) {
	// Validamos que la imagen no esta insertada en base de datos
	filter := bson.M{
		"name":      dtoInsertImage.Name,
		"owner":     dtoInsertImage.Owner,
		"extension": dtoInsertImage.Extension,
	}

	// Si la imagen ya existe, no la insertamos
	results, err := r.find(filter)
	if err != nil && err.Status != 404 {
		return nil, err
	}

	if err == nil && len(results) > 0 {
		return nil, exception.NewApiException(409, "La imagen ya existe")
	}

	image, errBuilder := imageBuilder.NewImageBuilder().
		FromImageUploadRequestDTO(dtoInsertImage).
		SetThumbnailId(thumbnailImageID).
		BuildNew()

	if errBuilder != nil {
		errorMessage := fmt.Sprintf("Error al construir la imagen: %s", errBuilder.Error())
		return nil, exception.NewApiException(404, errorMessage)
	}

	dto := imageDTO.FromImage(image)
	imageID, errInsert := r.mongoImage.InsertOne(context.Background(), dto)
	if errInsert != nil {
		return nil, exception.NewApiException(500, "Error al insertar el documento")
	}

	return &imageDTO.ImageUploadResponseDTO{
		Id:          imageID.InsertedID.(primitive.ObjectID).Hex(),
		ThumbnailId: thumbnailImageID,
		Name:        dto.Name,
		Extension:   dto.Extension,
		Size:        dto.Size,
	}, nil
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
