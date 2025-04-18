package controller

import (
	"go-gallery/src/commons/exception"
	"go-gallery/src/infrastructure/controller/handler"
	"go-gallery/src/infrastructure/dto"
	"go-gallery/src/service"

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

//	@Summary		Obtiene una imagen por su identificador
//	@Description	Obtiene una imagen específica del usuario según el identificador proporcionado
//	@Tags			image
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Identificador de la imagen"
//	@Security		CookieAuth
//	@Success		200	{object}	dto.DTOImage
//	@Failure		401	{object}	exception.ApiException	"Usuario no autenticado"
//	@Failure		403	{object}	exception.ApiException	"Los datos proporcionados no coinciden con el usuario autenticado"
//	@Failure		404	{object}	exception.ApiException	"Usuario/Imagen no encontrada"
//	@Failure		500	{object}	exception.ApiException	"Ha ocurrido un error inesperado"
//	@Router			/image/getImage/{id} [get]
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

//	@Summary		Persiste una imagen
//	@Description	Permite a un usuario autenticado persistir una imagen
//	@Tags			image
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file	true	"Archivo de imagen a subir (jpeg, jpg, png, webp)"
//	@Security		CookieAuth
//	@Success		200	{object}	dto.DTOImage			"Imagen subida correctamente"
//	@Failure		400	{object}	exception.ApiException	"Error al procesar la imagen"
//	@Failure		401	{object}	exception.ApiException	"Usuario no autenticado"
//	@Failure		403	{object}	exception.ApiException	"Los datos proporcionados no coinciden con el usuario autenticado"
//	@Failure		404	{object}	exception.ApiException	"Usuario/Imagen no encontrada"
//	@Failure		500	{object}	exception.ApiException	"Ha ocurrido un error inesperado"
//	@Router			/image/uploadImage [post]
func (c *ImageController) uploadImage(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(*dto.DTOClaimsJwt)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "Usuario no autenticado"))
	}

	_, errUser := c.serviceUser.FindAndCheckJWT(claims)
	if errUser != nil {
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	fileInput, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(exception.NewApiException(fiber.StatusNotFound, "Error al obtener la imagen del formulario"))
	}

	dtoInsertImage, errFile := handler.ProcessImageFile(fileInput, claims.Username)
	if errFile != nil {
		return ctx.Status(errFile.Status).JSON(errFile)
	}

	dto, errInsert := c.serviceImage.Insert(dtoInsertImage)
	if errInsert != nil {
		return ctx.Status(errInsert.Status).JSON(errInsert)
	}

	return ctx.Status(fiber.StatusOK).JSON(dto)
}

//	@Summary		Elimina una imagen
//	@Description	Borra una imagen específica del usuario autentificado
//	@Tags			image
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Identificador de la imagen"
//	@Security		CookieAuth
//	@Success		200	{object}	dto.DTOImage			"Imagen eliminada correctamente"
//	@Failure		401	{object}	exception.ApiException	"Usuario no autenticado"
//	@Failure		403	{object}	exception.ApiException	"Los datos proporcionados no coinciden con el usuario autenticado"
//	@Failure		404	{object}	exception.ApiException	"Usuario/Imagen no encontrada"
//	@Failure		500	{object}	exception.ApiException	"Ha ocurrido un error inesperado"
//	@Router			/image/deleteImage/{id} [delete]
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
