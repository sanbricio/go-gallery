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

LOCAL_MONGODB_URL=
MONGODB_DATABASE=

SECRET_KEY=
```

* `DOCKER_MONGO_INITDB_ROOT_USERNAME` y `DOCKER_MONGO_INITDB_ROOT_PASSWORD`: Credentials for the MongoDB root user.
* `DOCKER_ME_CONFIG_BASICAUTH_USERNAME` y `DOCKER_ME_CONFIG_BASICAUTH_PASSWORD`: Credentials to access the MongoDB Express website.
* `DOCKER_ME_CONFIG_MONGODB_URL`:  Connection URL to the MongoDB instance by docker
* `LOCAL_MONGODB_URL`:  Connection URL to the MongoDB local instance
* `MONGODB_DATABASE`:  Database name we are going to use
* `SECRET_KEY`: This key is used for JWT Authentification
