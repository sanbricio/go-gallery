package controller

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/controller/handler"
	"api-upload-photos/src/infrastructure/dto"
	"api-upload-photos/src/service"

	"github.com/gofiber/fiber/v2"
)

type ImageController struct {
	serviceImage *service.ImageService
	serviceUser  *service.UserService
}

func NewImageController(serviceImage *service.ImageService, serviceUser *service.UserService) *ImageController {
	return &ImageController{
		serviceImage: serviceImage,
		serviceUser:  serviceUser,
	}
}

func (c *ImageController) SetUpRoutes(router fiber.Router) {
	router.Get("/getImage/:id", c.getImage)
	router.Post("/uploadImage", c.uploadImage)
	router.Delete("/deleteImage/:id", c.deleteImage)
}

func (c *ImageController) getImage(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(*dto.DTOClaimsJwt)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "Usuario no autenticado"))
	}

	_, errUser := c.serviceUser.FindAndCheckJWT(claims)
	if errUser != nil {
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	dtoFindImage := &dto.DTOImage{
		IdImage: ctx.Params("id"),
		Owner:   claims.Username,
	}

	image, err := c.serviceImage.Find(dtoFindImage)
	if err != nil {
		return ctx.Status(err.Status).JSON(err)
	}

	return ctx.Status(fiber.StatusOK).JSON(image)
}

func (c *ImageController) uploadImage(ctx *fiber.Ctx) error {
	dtoUserJWT, ok := ctx.Locals("user").(*dto.DTOClaimsJwt)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "Usuario no autenticado"))
	}

	_, errUser := c.serviceUser.FindAndCheckJWT(dtoUserJWT)
	if errUser != nil {
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	fileInput, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(exception.NewApiException(fiber.StatusNotFound, "Error al obtener la imagen del formulario"))
	}

	dtoInsertImage, errFile := handler.ProcessImageFile(fileInput, dtoUserJWT.Username)
	if errFile != nil {
		return ctx.Status(errFile.Status).JSON(errFile)
	}

	dto, errInsert := c.serviceImage.Insert(dtoInsertImage)
	if errInsert != nil {
		return ctx.Status(errInsert.Status).JSON(errInsert)
	}

	return ctx.Status(fiber.StatusOK).JSON(dto)
}

func (c *ImageController) deleteImage(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(*dto.DTOClaimsJwt)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "Usuario no autenticado"))
	}

	_, errUser := c.serviceUser.FindAndCheckJWT(claims)
	if errUser != nil {
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	dtoFindImage := &dto.DTOImage{
		IdImage: ctx.Params("id"),
		Owner:   claims.Username,
	}
	image, err := c.serviceImage.Delete(dtoFindImage)
	if err != nil {
		return ctx.Status(err.Status).JSON(err)
	}

	return ctx.Status(fiber.StatusOK).JSON(image)
}
