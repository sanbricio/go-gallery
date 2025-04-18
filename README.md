# go-gallery

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
POSTGRESQL_HOST=localhost | postgres
POSTGRESQL_PORT=

MONGODB_URL_CONNECTION=mongodb://root:example@localhost:27017/ | mongodb://root:example@mongodb:27017/
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

- MongoDB Configuration:  
  - DOCKER_MONGO_INITDB_ROOT_USERNAME & DOCKER_MONGO_INITDB_ROOT_PASSWORD: Credentials for the MongoDB root user.  
  - DOCKER_ME_CONFIG_BASICAUTH_USERNAME & DOCKER_ME_CONFIG_BASICAUTH_PASSWORD: Credentials to access the MongoDB Express web interface.  
  - MONGODB_URL_CONNECTION: MongoDB connection string. Use mongodb://<user>:<password>@localhost:27017/ for local development, or mongodb://<user>:<password>@mongodb:27017/ when running with Docker Compose (where mongodb is the container name).
  - MONGODB_DATABASE: Name of the MongoDB database used by the application.  

- PostgreSQL Configuration:  
  - POSTGRESQL_USER: PostgreSQL username.  
  - POSTGRESQL_PASSWORD: PostgreSQL password.  
  - POSTGRESQL_DB: PostgreSQL database name (go-gallery).  
  - POSTGRESQL_HOST: Defines the PostgreSQL host. Use localhost for local development, or postgres if you're running the service inside Docker Compose (matching the container name).
  - POSTGRESQL_PORT: Port for PostgreSQL.  

- Email Sender Configuration:  
  - EMAIL_SENDER_HOST: Host for the email service (e.g., SMTP server).
  - EMAIL_SENDER_PORT: Port for sending emails. 
  - EMAIL_SENDER_USERNAME: Email address used to send emails.  
  - EMAIL_SENDER_PASSWORD: Password for the email account used to send emails.  

- Security & Authentication:  
  - JWT_SECRET: Secret key used for JWT authentication.  

- Application Configuration:  
  - GO_GALLERY_API_PORT: Port for the application.  
  - USER_REPOSITORY: Specifies the user repository implementation to use.  
  - IMAGE_REPOSITORY: Specifies the image repository implementation to use.  
  - EMAIL_SENDER_REPOSITORY: Specifies the email sender repository implementation to use.  

---

## 🛠️ Environment Setup

### 🔁 Docker Compose Configuration

When using Docker Compose, make sure to define the following variables in a .env file at the root of the project. Here is a recommended example:
```dotenv
DOCKER_MONGO_INITDB_ROOT_USERNAME=root
DOCKER_MONGO_INITDB_ROOT_PASSWORD=example
DOCKER_ME_CONFIG_BASICAUTH_USERNAME=admin
DOCKER_ME_CONFIG_BASICAUTH_PASSWORD=admin

POSTGRESQL_USER=postgres
POSTGRESQL_PASSWORD=postgres
POSTGRESQL_DB=go-gallery
POSTGRESQL_HOST=postgres
POSTGRESQL_PORT=5432

MONGODB_URL_CONNECTION=mongodb://root:example@mongodb:27017/
MONGODB_DATABASE=api-upload-images

EMAIL_SENDER_HOST=smtp.gmail.com
EMAIL_SENDER_PORT=587
EMAIL_SENDER_USERNAME=your_email@gmail.com
EMAIL_SENDER_PASSWORD=your_email_password

JWT_SECRET=your_super_secret_key
GO_GALLERY_API_PORT=3000
USER_REPOSITORY=UserPostgreSQLRepository
IMAGE_REPOSITORY=ImageMongoDBRepository
EMAIL_SENDER_REPOSITORY=EmailSenderGoMailRepository
```
> Ensure the service names in docker-compose.yml match postgres and mongodb so the containers can communicate properly.

---

### 💻 Local Development Configuration

If you're running the backend locally without Docker, you can use a .env file with the following configuration:
```dotenv
MONGODB_URL_CONNECTION=mongodb://root:example@localhost:27017/
MONGODB_DATABASE=api-upload-images

POSTGRESQL_USER=postgres
POSTGRESQL_PASSWORD=postgres
POSTGRESQL_DB=go-gallery
POSTGRESQL_HOST=localhost
POSTGRESQL_PORT=5432

EMAIL_SENDER_HOST=smtp.gmail.com
EMAIL_SENDER_PORT=587
EMAIL_SENDER_USERNAME=your_email@gmail.com
EMAIL_SENDER_PASSWORD=your_email_password

JWT_SECRET=your_super_secret_key
GO_GALLERY_API_PORT=3000
USER_REPOSITORY=UserPostgreSQLRepository
IMAGE_REPOSITORY=ImageMongoDBRepository
EMAIL_SENDER_REPOSITORY=EmailSenderGoMailRepository
```
> Make sure MongoDB and PostgreSQL are running locally.