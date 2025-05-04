package imageController

import (
	"go-gallery/src/commons/exception"
	imageService "go-gallery/src/service/image"
	userService "go-gallery/src/service/user"
	"strconv"

	imageDTO "go-gallery/src/infrastructure/dto/image"
	userDTO "go-gallery/src/infrastructure/dto/user"

	imageHandler "go-gallery/src/infrastructure/controller/image/handler"

	"github.com/gofiber/fiber/v2"
)

type ImageController struct {
	serviceImage *imageService.ImageService
	serviceUser  *userService.UserService
}

func NewImageController(serviceImage *imageService.ImageService, serviceUser *userService.UserService) *ImageController {
	return &ImageController{
		serviceImage: serviceImage,
		serviceUser:  serviceUser,
	}
}

func (c *ImageController) SetUpRoutes(router fiber.Router) {
	router.Get("/getImage/:id", c.getImage)
	router.Post("/uploadImage", c.uploadImage)
	router.Delete("/deleteImage/:id", c.deleteImage)

	// Thumbnails
	router.Get("/getThumbnailImages", c.getThumbnailImages)
}

// @Summary		Obtiene una imagen por su identificador
// @Description	Obtiene una imagen específica del usuario según el identificador proporcionado
// @Tags			image
// @Accept			json
// @Produce		json
// @Param			id	path	string	true	"Identificador de la imagen"
// @Security		CookieAuth
// @Success		200	{object}	imageDTO.ImageDTO
// @Failure		401	{object}	exception.ApiException	"Usuario no autenticado"
// @Failure		403	{object}	exception.ApiException	"Los datos proporcionados no coinciden con el usuario autenticado"
// @Failure		404	{object}	exception.ApiException	"Usuario/Imagen no encontrada"
// @Failure		500	{object}	exception.ApiException	"Ha ocurrido un error inesperado"
// @Router			/image/getImage/{id} [get]
func (c *ImageController) getImage(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(*userDTO.JwtClaimsDTO)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "Usuario no autenticado"))
	}

	_, errUser := c.serviceUser.FindAndCheckJWT(claims)
	if errUser != nil {
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	id := ctx.Params("id")
	dtoFindImage := &imageDTO.ImageDTO{
		Id:    &id,
		Owner: claims.Username,
	}

	image, err := c.serviceImage.Find(dtoFindImage)
	if err != nil {
		return ctx.Status(err.Status).JSON(err)
	}

	return ctx.Status(fiber.StatusOK).JSON(image)
}

// @Summary		Persiste una imagen
// @Description	Permite a un usuario autenticado persistir una imagen
// @Tags			image
// @Accept			multipart/form-data
// @Produce		json
// @Param			file	formData	file	true	"Archivo de imagen a subir (jpeg, jpg, png, webp)"
// @Security		CookieAuth
// @Success		200	{object}	imageDTO.ImageDTO		"Imagen subida correctamente"
// @Failure		400	{object}	exception.ApiException	"Error al procesar la imagen"
// @Failure		401	{object}	exception.ApiException	"Usuario no autenticado"
// @Failure		403	{object}	exception.ApiException	"Los datos proporcionados no coinciden con el usuario autenticado"
// @Failure		404	{object}	exception.ApiException	"Usuario/Imagen no encontrada"
// @Failure		409	{object}	exception.ApiException	"La imagen ya existe"
// @Failure		500	{object}	exception.ApiException	"Ha ocurrido un error inesperado"
// @Router			/image/uploadImage [post]
func (c *ImageController) uploadImage(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(*userDTO.JwtClaimsDTO)
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

	dtoInsertImage, errFile := imageHandler.ProcessImageFile(fileInput, claims.Username)
	if errFile != nil {
		return ctx.Status(errFile.Status).JSON(errFile)
	}

	dto, errInsert := c.serviceImage.Insert(dtoInsertImage)
	if errInsert != nil {
		return ctx.Status(errInsert.Status).JSON(errInsert)
	}

	return ctx.Status(fiber.StatusOK).JSON(dto)
}

// @Summary		Elimina una imagen
// @Description	Borra una imagen específica del usuario autentificado
// @Tags			image
// @Accept			json
// @Produce		json
// @Param			id	path	string	true	"Identificador de la imagen"
// @Security		CookieAuth
// @Success		200	{object}	imageDTO.ImageDTO		"Imagen eliminada correctamente"
// @Failure		401	{object}	exception.ApiException	"Usuario no autenticado"
// @Failure		403	{object}	exception.ApiException	"Los datos proporcionados no coinciden con el usuario autenticado"
// @Failure		404	{object}	exception.ApiException	"Usuario/Imagen no encontrada"
// @Failure		500	{object}	exception.ApiException	"Ha ocurrido un error inesperado"
// @Router			/image/deleteImage/{id} [delete]
func (c *ImageController) deleteImage(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(*userDTO.JwtClaimsDTO)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "Usuario no autenticado"))
	}

	_, errUser := c.serviceUser.FindAndCheckJWT(claims)
	if errUser != nil {
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	id := ctx.Params("id")
	dtoFindImage := &imageDTO.ImageDTO{
		Id:    &id,
		Owner: claims.Username,
	}
	image, err := c.serviceImage.Delete(dtoFindImage)
	if err != nil {
		return ctx.Status(err.Status).JSON(err)
	}

	return ctx.Status(fiber.StatusOK).JSON(image)
}

// @Summary		Listar imágenes en miniatura (thumbnails)
// @Description	Obtiene una lista paginada de imágenes en miniatura del usuario autenticado, usando paginación por cursor (lastId y pageSize).
// @Tags			thumbnail
// @Accept			json
// @Produce		json
// @Param			lastID		query	string	false	"Último ID recibido para la paginación"
// @Param			pageSize	query	int		false	"Cantidad de miniaturas a devolver (por defecto 10)"
// @Security		CookieAuth
// @Success		200	{object}	thumbnailImageDTO.ThumbnailImageCursorDTO	"Lista de miniaturas con el último id para poder realizar paginacione"
// @Failure		401	{object}	exception.ApiException						"Usuario no autenticado"
// @Failure		403	{object}	exception.ApiException						"Los datos proporcionados no coinciden con el usuario autenticado"
// @Failure		404	{object}	exception.ApiException						"No se encontraron thumbnails"
// @Failure		500	{object}	exception.ApiException						"Error inesperado"
// @Router			/image/getThumbnailImages [get]
func (c *ImageController) getThumbnailImages(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(*userDTO.JwtClaimsDTO)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "Usuario no autenticado"))
	}

	_, errUser := c.serviceUser.FindAndCheckJWT(claims)
	if errUser != nil {
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	// Obtener parámetros de consulta
	lastID := ctx.Query("lastID")
	pageSizeParam := ctx.Query("pageSize")

	// Validamos el pageSize (debe ser un entero positivo, default 10)
	pageSize := int64(10)
	if pageSizeParam != "" {
		if parsedPageSize, err := strconv.ParseInt(pageSizeParam, 10, 64); err == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
		}
	}

	thumbnails, errThumb := c.serviceImage.FindAllThumbnails(claims.Username, lastID, pageSize)
	if errThumb != nil {
		return ctx.Status(errThumb.Status).JSON(errThumb)
	}

	return ctx.Status(fiber.StatusOK).JSON(thumbnails)
}
