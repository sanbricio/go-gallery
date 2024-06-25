# API-Upload-Photos

This repository contains an API that allows uploading, retrieving, and deleting photos with persistence using either in-memory storage or MongoDB.

## Environment Variables Configuration

To deploy MongoDB and MongoDB Express interface, ensure the following environment variables are defined:

```dotenv
MONGO_INITDB_ROOT_USERNAME=
MONGO_INITDB_ROOT_PASSWORD=

ME_CONFIG_BASICAUTH_USERNAME=
ME_CONFIG_BASICAUTH_PASSWORD=
ME_CONFIG_MONGODB_URL=
```

* `MONGO_INITDB_ROOT_USERNAME` y `MONGO_INITDB_ROOT_PASSWORD`: Credentials for the MongoDB root user.
* `ME_CONFIG_BASICAUTH_USERNAME` y `ME_CONFIG_BASICAUTH_PASSWORD`: Credentials to access the MongoDB Express.
* `ME_CONFIG_MONGODB_URL`:  Connection URL to the MongoDB instance
