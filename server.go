package main

import (
	"api-upload-photos/src/infrastructure/controller"
	infrastructure "api-upload-photos/src/infrastructure/repository"
	"api-upload-photos/src/service"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	app := fiber.New()

	err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error cargando el archivo .env: %v", err)
    }

	mongoURL := os.Getenv("LOCAL_MONGODB_URL")
    databaseName := os.Getenv("MONGODB_DATABASE")

	//TODO Cambiar tipo de error a ConnectionException
	repository,errRepo:= infrastructure.NewRepositoryMongoDB(mongoURL,databaseName)
	if errRepo != nil {
		log.Fatal(errRepo.Message)
	}
	service := service.NewService(repository)
	controller := controller.NewController(app, service)
	controller.SetupRoutes()

	app.Listen(":3000")
}
