package imageRepository

import (
	"context"
	"fmt"
	"go-gallery/src/commons/exception"

	imageBuilder "go-gallery/src/domain/entities/builder/image"

	imageDTO "go-gallery/src/infrastructure/dto/image"
	log "go-gallery/src/infrastructure/logger"

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

var logger log.Logger

type ImageMongoDBRepository struct {
	mongoImage *mongo.Collection
}

func NewImageMongoDBRepository(args map[string]string) ImageRepository {
	urlConnection := args["MONGODB_URL_CONNECTION"]
	databaseName := args["MONGODB_DATABASE"]

	logger = log.Instance()

	db := connect(urlConnection, databaseName)

	repo := &ImageMongoDBRepository{
		mongoImage: db.Collection(IMAGE_COLLECTION),
	}

	logger.Info(fmt.Sprintf("Image repository initialized with connection to database '%s' and collection '%s'", databaseName, IMAGE_COLLECTION))
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

	logger.Info(fmt.Sprintf("Successfully connected to MongoDB with database '%s'", databaseName))
	return database.Database(databaseName)
}

func (r *ImageMongoDBRepository) Find(dtoFind *imageDTO.ImageDTO) (*imageDTO.ImageDTO, *exception.ApiException) {
	filter := bson.M{
		ID:    dtoFind.Id,
		OWNER: dtoFind.Owner,
	}

	logger.Info(fmt.Sprintf("Searching for image with filter: %+v", filter))

	result, err := r.find(filter)
	if err != nil {
		logger.Warning(fmt.Sprintf("Image not found with Id '%v' and Owner '%v'", dtoFind.Id, dtoFind.Owner))
		return nil, err
	}

	logger.Info(fmt.Sprintf("Image found: %+v", result[0]))
	return &result[0], nil
}

func (r *ImageMongoDBRepository) find(filter bson.M) ([]imageDTO.ImageDTO, *exception.ApiException) {
	cursor, err := r.mongoImage.Find(context.Background(), filter)
	if err != nil {
		logger.Error(fmt.Sprintf("Error searching for images with filter: %+v - %s", filter, err.Error()))
		return nil, exception.NewApiException(500, "Error searching for images")
	}
	defer cursor.Close(context.Background())

	var results []imageDTO.ImageDTO
	for cursor.Next(context.Background()) {
		var image imageDTO.ImageDTO
		if err := cursor.Decode(&image); err != nil {
			logger.Error(fmt.Sprintf("Error decoding image: %s", err.Error()))
			return nil, exception.NewApiException(500, "Error decoding images")
		}
		results = append(results, image)
	}

	if len(results) == 0 {
		logger.Warning(fmt.Sprintf("No results found with filter: %+v", filter))
		return nil, exception.NewApiException(404, "Image not found")
	}

	logger.Info(fmt.Sprintf("Found %d images", len(results)))
	return results, nil
}

func (r *ImageMongoDBRepository) Insert(dtoInsertImage *imageDTO.ImageUploadRequestDTO, thumbnailImageID string) (*imageDTO.ImageUploadResponseDTO, *exception.ApiException) {
	filter := bson.M{
		"name":      dtoInsertImage.Name,
		"owner":     dtoInsertImage.Owner,
		"extension": dtoInsertImage.Extension,
	}

	logger.Info(fmt.Sprintf("Checking if the image already exists with filter: %+v", filter))

	results, err := r.find(filter)
	if err != nil && err.Status != 404 {
		logger.Error(fmt.Sprintf("Error checking if the image already exists: %s", err.Message))
		return nil, err
	}

	if err == nil && len(results) > 0 {
		logger.Warning(fmt.Sprintf("Image already exists for Owner '%s' with name '%s'", dtoInsertImage.Owner, dtoInsertImage.Name))
		return nil, exception.NewApiException(409, "Image already exists")
	}

	logger.Info(fmt.Sprintf("Building image entity for owner '%s' with name '%s'", dtoInsertImage.Owner, dtoInsertImage.Name))

	image, errBuilder := imageBuilder.NewImageBuilder().
		FromImageUploadRequestDTO(dtoInsertImage).
		SetThumbnailId(thumbnailImageID).
		BuildNew()

	if errBuilder != nil {
		errorMessage := fmt.Sprintf("Error building image: %s", errBuilder.Error())
		logger.Error(errorMessage)
		return nil, exception.NewApiException(500, errorMessage)
	}

	dto := imageDTO.FromImage(image)

	logger.Info(fmt.Sprintf("Inserting new image into the database: %+v", dto))

	imageID, errInsert := r.mongoImage.InsertOne(context.Background(), dto)
	if errInsert != nil {
		logger.Error(fmt.Sprintf("Error inserting image: %s", errInsert.Error()))
		return nil, exception.NewApiException(500, "Error inserting the document")
	}

	logger.Info(fmt.Sprintf("Image successfully inserted with ID: %s", imageID.InsertedID.(primitive.ObjectID).Hex()))

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

	logger.Info(fmt.Sprintf("Attempting to delete image with filter: %+v", filter))

	foundImages, err := r.find(filter)
	if err != nil {
		logger.Warning(fmt.Sprintf("Image not found for deletion with Id '%v' and Owner '%v'", dto.Id, dto.Owner))
		return nil, err
	}

	_, errDelete := r.mongoImage.DeleteOne(context.Background(), filter)
	if errDelete != nil {
		logger.Error(fmt.Sprintf("Error deleting image: %s", errDelete.Error()))
		return nil, exception.NewApiException(500, "Error deleting the image")
	}

	logger.Info(fmt.Sprintf("Image successfully deleted: %+v", foundImages[0]))
	return &foundImages[0], nil
}
