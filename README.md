# GoGallery Backend

This repository contains an API that enables users to upload, retrieve, and delete photos. Additionally, it includes authentication functionality to ensure secure access

## Environment Variables Configuration

To deploy GoGallery, make sure to define the following environment variables:

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
* `GO_GALLERY_API_PORT`: This key is used for select a specified port for the application
* `USER_REPOSITORY`: This key is used for choose a implementation of user repository
* `IMAGE_REPOSITORY`: This key is used for choose a implementation of image repository
