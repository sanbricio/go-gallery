package thumbnailImageRepository

import (
	"context"
	"fmt"
	"go-gallery/src/commons/constants"
	"go-gallery/src/commons/exception"
	utilsImage "go-gallery/src/commons/utils/image"
	thumbnailImageBuilder "go-gallery/src/domain/entities/builder/image/thumbnailImage"
	"strings"

	imageDTO "go-gallery/src/infrastructure/dto/image"
	thumbnailImageDTO "go-gallery/src/infrastructure/dto/image/thumbnailImage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const ThumnbailImageMongoDBRepositoryKey = "ThumnbailImageMongoDBRepository"

const (
	THUMBNAIL_IMAGE_COLLECTION string = "ThumbnailImage"
	ID                         string = "id"
	OWNER                      string = "owner"
)

type ThumbnailImageMongoDBRepository struct {
	mongoThumbnailImage *mongo.Collection
}

func NewThumbnailImageMongoDBRepository(args map[string]string) ThumbnailImageRepository {
	urlConnection := args["MONGODB_URL_CONNECTION"]
	databaseName := args["MONGODB_DATABASE"]

	db := connect(urlConnection, databaseName)

	repo := &ThumbnailImageMongoDBRepository{
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

func (r *ThumbnailImageMongoDBRepository) find(filter bson.M) ([]thumbnailImageDTO.ThumbnailImageDTO, *exception.ApiException) {
	cursor, err := r.mongoThumbnailImage.Find(context.Background(), filter)
	if err != nil {
		return nil, exception.NewApiException(500, "Error al buscar las miniaturas")
	}
	defer cursor.Close(context.Background())

	var results []thumbnailImageDTO.ThumbnailImageDTO
	for cursor.Next(context.Background()) {
		var thumbnail thumbnailImageDTO.ThumbnailImageDTO
		if err := cursor.Decode(&thumbnail); err != nil {
			return nil, exception.NewApiException(500, "Error al decodificar el las miniaturas")
		}
		results = append(results, thumbnail)
	}

	if len(results) == 0 {
		return nil, exception.NewApiException(404, "Miniatura no encontrada")
	}

	return results, nil
}

func (r *ThumbnailImageMongoDBRepository) Insert(dto *imageDTO.ImageUploadRequestDTO) (string, *exception.ApiException) {
	// Validamos que la miniatura no esta insertada en base de datos
	filter := bson.M{
		"name":      strings.TrimSpace(dto.Name),
		"owner":     strings.TrimSpace(dto.Owner),
		"extension": strings.TrimSpace(dto.Extension),
	}

	// Si la miniatura ya existe, no la insertamos
	results, err := r.find(filter)
	if err != nil && err.Status != 404 {
		return "", err
	}

	if err == nil && len(results) > 0 {
		return "", exception.NewApiException(409, "La miniatura ya existe")
	}

	resizedBytes, errResize := utilsImage.ResizeImage(dto.RawContentFile, constants.THUMBNAIL_WIDTH, constants.THUMBNAIL_HEIGHT)
	if errResize != nil {
		errorMessage := fmt.Sprintf("Error al generar miniatura: %s", errResize.Error())
		return "", exception.NewApiException(500, errorMessage)
	}

	thumbnailImage, errBuilder := thumbnailImageBuilder.NewThumbnailImageBuilder().
		FromImageUploadRequestDTO(dto).
		SetContentFile(utilsImage.EncondeImageToBase64(resizedBytes)).
		BuildNew()

	if errBuilder != nil {
		errorMessage := fmt.Sprintf("Error al construir la miniatura: %s", errBuilder.Error())
		return "", exception.NewApiException(404, errorMessage)
	}

	dtoThumbnailImage := thumbnailImageDTO.FromThumbnailImage(thumbnailImage)

	thumbnailId, errInsert := r.mongoThumbnailImage.InsertOne(context.Background(), dtoThumbnailImage)
	if errInsert != nil {
		return "", exception.NewApiException(500, "Error al insertar el documento")
	}
	// Convertimos el ID de la miniatura a string
	return thumbnailId.InsertedID.(primitive.ObjectID).Hex(), nil
}
