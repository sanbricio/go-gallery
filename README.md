# API-Upload-Photos

This repository contains an API that allows uploading, retrieving, and deleting photos with persistence using either in-memory storage or MongoDB.

## Environment Variables Configuration

To deploy MongoDB and MongoDB Express interface, ensure the following environment variables are defined:

```dotenv
DOCKER_MONGO_INITDB_ROOT_USERNAME=
DOCKER_MONGO_INITDB_ROOT_PASSWORD=
DOCKER_ME_CONFIG_BASICAUTH_USERNAME=
DOCKER_ME_CONFIG_BASICAUTH_PASSWORD=
DOCKER_ME_CONFIG_MONGODB_URL=

MONGODB_URL_CONNECTION=
MONGODB_DATABASE=

JWT_SECRET=

UPLOAD_PHOTOS_API_USER_REPOSITORY=UserMongoDBRepository
UPLOAD_PHOTOS_API_IMAGE_REPOSITORY=ImageMongoDBRepository
```

* `DOCKER_MONGO_INITDB_ROOT_USERNAME` y `DOCKER_MONGO_INITDB_ROOT_PASSWORD`: Credentials for the MongoDB root user.
* `DOCKER_ME_CONFIG_BASICAUTH_USERNAME` y `DOCKER_ME_CONFIG_BASICAUTH_PASSWORD`: Credentials to access the MongoDB Express website.
* `DOCKER_ME_CONFIG_MONGODB_URL`:  Connection URL to the MongoDB instance by docker
* `MONGODB_URL_CONNECTION`:  Connection URL to the MongoDB local instance
* `MONGODB_DATABASE`:  Database name we are going to use
* `JWT_SECRET`: This key is used for JWT Authentification
* `UPLOAD_PHOTOS_API_PORT`: This key is used for select a specified port for the application
* `UPLOAD_PHOTOS_API_USER_REPOSITORY`: This key is used for choose a implementation of user repository
* `UPLOAD_PHOTOS_API_IMAGE_REPOSITORY`: This key is used for choose a implementation of image repository
