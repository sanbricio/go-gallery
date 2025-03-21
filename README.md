
# GoGallery Backend

This repository contains an API that enables users to upload, retrieve, and delete photos. Additionally, it includes authentication functionality to ensure secure access.

## Environment Variables Configuration

To deploy GoGallery, make sure to define the following environment variables:

```dotenv
DOCKER_MONGO_INITDB_ROOT_USERNAME=
DOCKER_MONGO_INITDB_ROOT_PASSWORD=
DOCKER_ME_CONFIG_BASICAUTH_USERNAME=
DOCKER_ME_CONFIG_BASICAUTH_PASSWORD=

POSTGRESQL_USER=
POSTGRESQL_PASSWORD=
POSTGRESQL_DB=go-gallery
POSTGRESQL_HOST=
POSTGRESQL_PORT=

MONGODB_URL_CONNECTION=mongodb://root:example@localhost:27017/
MONGODB_DATABASE=api-upload-images

JWT_SECRET=

GO_GALLERY_API_PORT=3000
USER_REPOSITORY=UserPostgreSQLRepository
IMAGE_REPOSITORY=ImageMongoDBRepository
EMAIL_SENDER_REPOSITORY=EmailSenderGoMailRepository

EMAIL_SENDER_HOST=smtp.gmail.com
EMAIL_SENDER_PORT=587
EMAIL_SENDER_USERNAME=
EMAIL_SENDER_PASSWORD=
```

### Description of Environment Variables

- **MongoDB Configuration:**  
  - `DOCKER_MONGO_INITDB_ROOT_USERNAME` & `DOCKER_MONGO_INITDB_ROOT_PASSWORD`: Credentials for the MongoDB root user.  
  - `DOCKER_ME_CONFIG_BASICAUTH_USERNAME` & `DOCKER_ME_CONFIG_BASICAUTH_PASSWORD`: Credentials to access the MongoDB Express web interface.  
  - `MONGODB_URL_CONNECTION`: Connection URL for the local MongoDB instance.  
  - `MONGODB_DATABASE`: Name of the MongoDB database used by the application.  

- **PostgreSQL Configuration:**  
  - `POSTGRESQL_USER`: PostgreSQL username.  
  - `POSTGRESQL_PASSWORD`: PostgreSQL password.  
  - `POSTGRESQL_DB`: PostgreSQL database name (`go-gallery`).  
  - `POSTGRESQL_HOST`: Host for PostgreSQL.  
  - `POSTGRESQL_PORT`: Port for PostgreSQL.  

- **Email Sender Configuration:**  
  - `EMAIL_SENDER_HOST`: Host for the email service (e.g., SMTP server).
  - `EMAIL_SENDER_PORT`: Port for sending emails. 
  - `EMAIL_SENDER_USERNAME`: Email address used to send emails..  
  - `EMAIL_SENDER_PASSWORD`: Password for the email account used to send emails.  

- **Security & Authentication:**  
  - `JWT_SECRET`: Secret key used for JWT authentication.  

- **Application Configuration:**  
  - `GO_GALLERY_API_PORT`: Port for the application.  
  - `USER_REPOSITORY`: Specifies the user repository implementation to use.  
  - `IMAGE_REPOSITORY`: Specifies the image repository implementation to use.  
  - `EMAIL_SENDER_REPOSITORY`: Specifies the email sender repository implementation to use.  

