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
	log "go-gallery/src/infrastructure/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const ThumnbailImageMongoDBRepositoryKey = "ThumnbailImageMongoDBRepository"

const (
	THUMBNAIL_IMAGE_COLLECTION string = "ThumbnailImage"
	ID                         string = "_id"
	OWNER                      string = "owner"
	SORT                       int    = -1 // Ordenado de manera descendente (mas reciente primero)
)

var logger log.Logger

type ThumbnailImageMongoDBRepository struct {
	mongoThumbnailImage *mongo.Collection
}

func NewThumbnailImageMongoDBRepository(args map[string]string) ThumbnailImageRepository {
	urlConnection := args["MONGODB_URL_CONNECTION"]
	databaseName := args["MONGODB_DATABASE"]

	logger = log.Instance()
	logger.Info("Initializing ThumbnailImageMongoDBRepository with MongoDB URL: " + urlConnection)

	db := connect(urlConnection, databaseName)

	repo := &ThumbnailImageMongoDBRepository{
		mongoThumbnailImage: db.Collection(THUMBNAIL_IMAGE_COLLECTION),
	}

	logger.Info("ThumbnailImageMongoDBRepository successfully initialized")

	return repo
}

func connect(urlConnection string, databaseName string) *mongo.Database {
	database, err := mongo.Connect(context.Background(), options.Client().ApplyURI(urlConnection))
	if err != nil {
		panicMessage := fmt.Sprintf("Could not connect to MongoDB: %s", err.Error())
		logger.Panic(panicMessage)
		panic(panicMessage)
	}

	err = database.Ping(context.Background(), readpref.Primary())
	if err != nil {
		panicMessage := fmt.Sprintf("Could not ping MongoDB: %s", err.Error())
		logger.Panic(panicMessage)
		panic(panicMessage)
	}

	logger.Info("Successfully connected to MongoDB")

	return database.Database(databaseName)
}

func (r *ThumbnailImageMongoDBRepository) FindAll(owner, lastID string, pageSize int64) (*thumbnailImageDTO.ThumbnailImageCursorDTO, *exception.ApiException) {

	filter := bson.M{
		"owner": strings.TrimSpace(owner),
	}

	logger.Info(fmt.Sprintf("Cursor-based search: filter=%+v, lastID=%s, pageSize=%d", filter, lastID, pageSize))

	if lastID != "" {
		lastObjectID, err := primitive.ObjectIDFromHex(lastID)
		if err != nil {
			logger.Error(fmt.Sprintf("Invalid lastID: %s", lastID))
			return nil, exception.NewApiException(400, "Invalid last ID")
		}
		filter[ID] = bson.M{"$lt": lastObjectID}
	}

	findOptions := options.Find()
	findOptions.SetLimit(pageSize)
	findOptions.SetSort(bson.D{{Key: ID, Value: SORT}})

	dto, err := r.find(filter, findOptions)
	if err != nil {
		return nil, err
	}

	return &thumbnailImageDTO.ThumbnailImageCursorDTO{
		Thumbnails: dto,
		LastID:     *dto[len(dto)-1].Id,
	}, nil
}

func (r *ThumbnailImageMongoDBRepository) find(filter bson.M, findOptions *options.FindOptions) ([]thumbnailImageDTO.ThumbnailImageDTO, *exception.ApiException) {
	logger.Info(fmt.Sprintf("Searching for thumbnails with filter: %+v and options: %+v", filter, findOptions))
	cursor, err := r.mongoThumbnailImage.Find(context.Background(), filter, findOptions)
	if err != nil {
		logger.Error(fmt.Sprintf("Error searching for thumbnails: %s", err.Error()))
		return nil, exception.NewApiException(500, "Error searching for thumbnails")
	}
	defer cursor.Close(context.Background())

	var results []thumbnailImageDTO.ThumbnailImageDTO
	for cursor.Next(context.Background()) {
		var thumbnail thumbnailImageDTO.ThumbnailImageDTO
		if err := cursor.Decode(&thumbnail); err != nil {
			logger.Error(fmt.Sprintf("Error decoding thumbnail: %s", err.Error()))
			return nil, exception.NewApiException(500, "Error decoding thumbnails")
		}
		results = append(results, thumbnail)
	}

	if len(results) == 0 {
		logger.Warning("No thumbnails found")
		return nil, exception.NewApiException(404, "Thumbnail not found")
	}

	logger.Info(fmt.Sprintf("Thumbnails found: %v", len(results)))
	return results, nil
}

func (r *ThumbnailImageMongoDBRepository) Insert(dto *imageDTO.ImageUploadRequestDTO) (string, *exception.ApiException) {
	// Validamos que la miniatura no esta insertada en base de datos
	filter := bson.M{
		"name":      strings.TrimSpace(dto.Name),
		"owner":     strings.TrimSpace(dto.Owner),
		"extension": strings.TrimSpace(dto.Extension),
	}

	logger.Info(fmt.Sprintf("Attempting to insert thumbnail for image: Name=%s, Owner=%s", dto.Name, dto.Owner))
	// Si la miniatura ya existe, no la insertamos
	results, err := r.find(filter, nil)
	if err != nil && err.Status != 404 {
		return "", err
	}

	if err == nil && len(results) > 0 {
		logger.Warning("Thumbnail already exists and will not be inserted")
		return "", exception.NewApiException(409, "Thumbnail already exists")
	}

	resizedBytes, errResize := utilsImage.ResizeImage(dto.RawContentFile, constants.THUMBNAIL_WIDTH, constants.THUMBNAIL_HEIGHT)
	if errResize != nil {
		errorMessage := fmt.Sprintf("Error generating thumbnail: %s", errResize.Error())
		logger.Error(errorMessage)
		return "", exception.NewApiException(500, errorMessage)
	}

	sizeInBytes := len(resizedBytes)
	size := utilsImage.HumanizeBytes(uint64(sizeInBytes))

	thumbnailImage, errBuilder := thumbnailImageBuilder.NewThumbnailImageBuilder().
		FromImageUploadRequestDTO(dto).
		SetContentFile(utilsImage.EncondeImageToBase64(resizedBytes)).
		SetSize(size).
		BuildNew()

	if errBuilder != nil {
		errorMessage := fmt.Sprintf("Error building thumbnail: %s", errBuilder.Error())
		logger.Error(errorMessage)
		return "", exception.NewApiException(404, errorMessage)
	}

	dtoThumbnailImage := thumbnailImageDTO.FromThumbnailImage(thumbnailImage)

	thumbnailId, errInsert := r.mongoThumbnailImage.InsertOne(context.Background(), dtoThumbnailImage)
	if errInsert != nil {
		logger.Error(fmt.Sprintf("Error inserting thumbnail: %s", errInsert.Error()))
		return "", exception.NewApiException(500, "Error inserting document")
	}

	idHex := thumbnailId.InsertedID.(primitive.ObjectID).Hex()
	logger.Info(fmt.Sprintf("Thumbnail successfully inserted with ID: %s", idHex))
	return idHex, nil
}
