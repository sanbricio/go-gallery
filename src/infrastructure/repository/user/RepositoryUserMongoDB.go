package repository

import (
	"api-upload-photos/src/commons/exception"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type RepositoryUserMongoDB struct {
	client *mongo.Database
}

func NewRepositoryMongoDB(urlConnection string, databaseName string) (*RepositoryUserMongoDB, *exception.ConnectionException) {
	db, err := connect(urlConnection, databaseName)
	if err != nil {
		return nil, err
	}

	repo := &RepositoryUserMongoDB{
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
