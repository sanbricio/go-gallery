// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Support GoGallery",
            "email": "gogalleryteam@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/delete": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Elimina la cuenta de usuario tras verificar el código enviado",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Confirmar eliminación de cuenta",
                "parameters": [
                    {
                        "description": "Datos para confirmar eliminación",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.DTODeleteUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Se han eliminado los datos del usuario correctamente",
                        "schema": {
                            "$ref": "#/definitions/dto.DTOMessageResponse"
                        }
                    },
                    "400": {
                        "description": "Solicitud incorrecta",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "401": {
                        "description": "Usuario no autenticado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "403": {
                        "description": "Los datos proporcionados no coinciden con el usuario autenticado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "404": {
                        "description": "Usuario no encontrado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "500": {
                        "description": "Ha ocurrido un error inesperado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "Autentica un usuario y genera un token JWT para guardarlo en una cookie",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Iniciar sesión",
                "parameters": [
                    {
                        "description": "Datos de autenticación",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.DTOLoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Se ha iniciado sesion correctamente",
                        "schema": {
                            "$ref": "#/definitions/dto.DTOLoginResponse"
                        },
                        "headers": {
                            "Set-Cookie": {
                                "type": "string",
                                "description": "Authorization=auth_token; HttpOnly; Secure"
                            }
                        }
                    },
                    "400": {
                        "description": "Contraseña incorrecta",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "401": {
                        "description": "No autorizado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "404": {
                        "description": "Usuario no encontrado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "500": {
                        "description": "Ha ocurrido un error inesperado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "post": {
                "security": [
                    {
                        "CookieAuth": []
                    }
                ],
                "description": "Cierra la sesión del usuario autenticado, elimina la cookie auth_token",
                "tags": [
                    "auth"
                ],
                "summary": "Cerrar sesión",
                "responses": {
                    "200": {
                        "description": "Se ha cerrado sesión correctamente",
                        "schema": {
                            "$ref": "#/definitions/dto.DTOMessageResponse"
                        }
                    },
                    "401": {
                        "description": "Usuario no autenticado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "403": {
                        "description": "Los datos proporcionados no coinciden con el usuario autenticado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "404": {
                        "description": "Usuario no encontrado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "500": {
                        "description": "Ha ocurrido un error inesperado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Registra un nuevo usuario en el sistema",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Registro de un nuevo usuario",
                "parameters": [
                    {
                        "description": "Datos de registro",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.DTOUser"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Usuario creado",
                        "schema": {
                            "$ref": "#/definitions/dto.DTORegisterResponse"
                        }
                    },
                    "400": {
                        "description": "Solicitud incorrecta",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "500": {
                        "description": "Ha ocurrido un error inesperado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    }
                }
            }
        },
        "/auth/request-delete": {
            "post": {
                "security": [
                    {
                        "CookieAuth": []
                    }
                ],
                "description": "Envía un código de verificación al correo para eliminar la cuenta",
                "tags": [
                    "auth"
                ],
                "summary": "Solicitar eliminación de cuenta",
                "responses": {
                    "200": {
                        "description": "Se ha enviado un código de confirmación al correo electrónico",
                        "schema": {
                            "$ref": "#/definitions/dto.DTOMessageResponse"
                        }
                    },
                    "401": {
                        "description": "Usuario no autenticado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "403": {
                        "description": "Los datos proporcionados no coinciden con el usuario autenticado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "404": {
                        "description": "Usuario no encontrado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "500": {
                        "description": "Ha ocurrido un error inesperado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    }
                }
            }
        },
        "/auth/update": {
            "put": {
                "security": [
                    {
                        "CookieAuth": []
                    }
                ],
                "description": "Actualiza los datos de un usuario autenticado",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Actualizar usuario",
                "parameters": [
                    {
                        "description": "Datos de actualización",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.DTOUpdateUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Se han actualizado los datos del usuario correctamente.",
                        "schema": {
                            "$ref": "#/definitions/dto.DTOMessageResponse"
                        }
                    },
                    "400": {
                        "description": "Solicitud incorrecta",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "401": {
                        "description": "Usuario no autenticado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "403": {
                        "description": "Los datos proporcionados no coinciden con el usuario autenticado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "404": {
                        "description": "Usuario no encontrado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "500": {
                        "description": "Ha ocurrido un error inesperado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    }
                }
            }
        },
        "/docs/definition/swagger.json": {
            "get": {
                "description": "Retorna la definición de la API(OpenAPI) en formato JSON",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "docs"
                ],
                "summary": "Obtiene la documentación de la API en formato JSON",
                "responses": {
                    "200": {
                        "description": "Archivo JSON cargado correctamente"
                    },
                    "500": {
                        "description": "Error al cargar el archivo JSON"
                    }
                }
            }
        },
        "/docs/definition/swagger.yml": {
            "get": {
                "description": "Retorna la definición de la API(OpenAPI) en formato YAML",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "docs"
                ],
                "summary": "Obtiene la documentación de la API en formato YAML",
                "responses": {
                    "200": {
                        "description": "Archivo YAML cargado correctamente"
                    },
                    "500": {
                        "description": "Error al cargar el archivo YAML"
                    }
                }
            }
        },
        "/image/deleteImage/{id}": {
            "delete": {
                "security": [
                    {
                        "CookieAuth": []
                    }
                ],
                "description": "Borra una imagen específica del usuario autentificado",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "image"
                ],
                "summary": "Elimina una imagen",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Identificador de la imagen",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Imagen eliminada correctamente",
                        "schema": {
                            "$ref": "#/definitions/dto.DTOImage"
                        }
                    },
                    "401": {
                        "description": "Usuario no autenticado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "403": {
                        "description": "Los datos proporcionados no coinciden con el usuario autenticado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "404": {
                        "description": "Usuario/Imagen no encontrada",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "500": {
                        "description": "Ha ocurrido un error inesperado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    }
                }
            }
        },
        "/image/getImage/{id}": {
            "get": {
                "security": [
                    {
                        "CookieAuth": []
                    }
                ],
                "description": "Obtiene una imagen específica del usuario según el identificador proporcionado",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "image"
                ],
                "summary": "Obtiene una imagen por su identificador",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Identificador de la imagen",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.DTOImage"
                        }
                    },
                    "401": {
                        "description": "Usuario no autenticado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "403": {
                        "description": "Los datos proporcionados no coinciden con el usuario autenticado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "404": {
                        "description": "Usuario/Imagen no encontrada",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "500": {
                        "description": "Ha ocurrido un error inesperado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    }
                }
            }
        },
        "/image/uploadImage": {
            "post": {
                "security": [
                    {
                        "CookieAuth": []
                    }
                ],
                "description": "Permite a un usuario autenticado persistir una imagen",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "image"
                ],
                "summary": "Persiste una imagen",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Archivo de imagen a subir (jpeg, jpg, png, webp)",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Imagen subida correctamente",
                        "schema": {
                            "$ref": "#/definitions/dto.DTOImage"
                        }
                    },
                    "400": {
                        "description": "Error al procesar la imagen",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "401": {
                        "description": "Usuario no autenticado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "403": {
                        "description": "Los datos proporcionados no coinciden con el usuario autenticado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "404": {
                        "description": "Usuario/Imagen no encontrada",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    },
                    "500": {
                        "description": "Ha ocurrido un error inesperado",
                        "schema": {
                            "$ref": "#/definitions/exception.ApiException"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.DTODeleteUser": {
            "description": "Datos necesarios para proceder con la eliminación del usuario",
            "type": "object",
            "properties": {
                "code": {
                    "description": "Código de verificación para la eliminación\nexample \"123456\"",
                    "type": "string",
                    "example": "123456"
                },
                "password": {
                    "description": "Contraseña del usuario para confirmar la eliminación\nexample \"MiContraseñaSegura\"",
                    "type": "string",
                    "example": "MiContraseñaSegura"
                }
            }
        },
        "dto.DTOImage": {
            "description": "Contiene la información de una imagen, incluyendo su identificador, nombre, extensión, contenido en base64 y propietario (usuario)",
            "type": "object",
            "properties": {
                "content_file": {
                    "description": "Contenido de la imagen en base64\nExample: /9j/4AAQSkZJRgABAQEAAAAAAAD...",
                    "type": "string",
                    "example": "/9j/4AAQSkZJRgABAQEAAAAAAAD."
                },
                "extension": {
                    "description": "Extensión del archivo de imagen\nExample: jpg",
                    "type": "string",
                    "example": ".jpeg"
                },
                "id_image": {
                    "description": "ID único de la imagen\nExample: 64a1f8b8e4b0c10d3c5b2e75",
                    "type": "string",
                    "example": "64a1f8b8e4b0c10d3c5b2e75"
                },
                "name": {
                    "description": "Nombre del archivo de la imagen\nExample: foto_perfil",
                    "type": "string",
                    "example": "prueba"
                },
                "owner": {
                    "description": "Usuario propietario de la imagen\nExample: usuario123",
                    "type": "string",
                    "example": "usuario123"
                },
                "size": {
                    "description": "Tamaño de la imagen en bytes\nExample: 204800",
                    "type": "string",
                    "example": "2.3 kB"
                }
            }
        },
        "dto.DTOLoginRequest": {
            "description": "Datos requeridos para realizar autentificación del usuario",
            "type": "object",
            "properties": {
                "password": {
                    "description": "Contraseña del usuario\nexample \"MiContraseñaSegura.\"",
                    "type": "string",
                    "example": "MiContraseñaSegura."
                },
                "username": {
                    "description": "Nombre de usuario\nexample \"usuario\"",
                    "type": "string",
                    "example": "usuario"
                }
            }
        },
        "dto.DTOLoginResponse": {
            "description": "Respuesta al iniciar sesión correctamente",
            "type": "object",
            "properties": {
                "email": {
                    "description": "Correo electrónico del usuario\nexample \"usuario@example.com\"",
                    "type": "string",
                    "example": "usuario@example.com"
                },
                "message": {
                    "description": "Mensaje de confirmación\nexample \"Se ha iniciado sesión correctamente\"",
                    "type": "string",
                    "example": "Se ha iniciado sesión correctamente"
                },
                "name": {
                    "description": "Nombre del usuario\nexample \"Juan Pérez\"",
                    "type": "string",
                    "example": "Juan Pérez"
                },
                "username": {
                    "description": "Nombre de usuario autenticado\nexample \"usuario123\"",
                    "type": "string",
                    "example": "usuario123"
                }
            }
        },
        "dto.DTOMessageResponse": {
            "description": "Respuesta con un mensaje para informar al usuario que ha ocurrido",
            "type": "object",
            "properties": {
                "message": {
                    "description": "Mensaje de respuesta\nexample \"Operación realizada con éxito\"",
                    "type": "string",
                    "example": "Ha funcionado correctamente"
                }
            }
        },
        "dto.DTORegisterResponse": {
            "description": "Respuesta generada después de crear un nuevo usuario",
            "type": "object",
            "properties": {
                "firstname": {
                    "description": "Nombre\nexample \"Juan\"",
                    "type": "string",
                    "example": "Juan"
                },
                "message": {
                    "description": "Mensaje de confirmación\nexample \"Se ha creado el usuario correctamente\"",
                    "type": "string",
                    "example": "Se ha creado el usuario correctamente"
                },
                "username": {
                    "description": "Nombre de usuario\nexample \"usuario123\"",
                    "type": "string",
                    "example": "usuario123"
                }
            }
        },
        "dto.DTOUpdateUser": {
            "description": "Datos que pueden ser actualizados del usuario existente",
            "type": "object",
            "properties": {
                "email": {
                    "description": "Correo electrónico\nexample \"nuevo.email@example.com\"",
                    "type": "string",
                    "example": "nuevo.email@example.com"
                },
                "firstname": {
                    "description": "Nombre\nexample \"Carlos\"",
                    "type": "string",
                    "example": "Carlos"
                },
                "lastname": {
                    "description": "Apellido\nexample \"Gómez\"",
                    "type": "string",
                    "example": "Gómez"
                },
                "password": {
                    "description": "Contraseña\nexample \"NuevaContraseñaSegura.\"",
                    "type": "string",
                    "example": "NuevaContraseñaSegura."
                }
            }
        },
        "dto.DTOUser": {
            "description": "Estructura que define los datos del usuario",
            "type": "object",
            "properties": {
                "email": {
                    "description": "Correo electrónico\nexample \"usuario@example.com\"",
                    "type": "string",
                    "example": "usuario@example.com"
                },
                "firstname": {
                    "description": "Nombre\nexample \"Juan\"",
                    "type": "string",
                    "example": "Juan"
                },
                "lastname": {
                    "description": "Apellido\nexample \"Pérez\"",
                    "type": "string",
                    "example": "Pérez"
                },
                "password": {
                    "description": "Contraseña\nexample \"MiContraseñaSegura.\"",
                    "type": "string",
                    "example": "MiContraseñaSegura."
                },
                "username": {
                    "description": "Nombre de usuario\nexample \"usuario123\"",
                    "type": "string",
                    "example": "usuario123"
                }
            }
        },
        "exception.ApiException": {
            "description": "Estructura para manejar excepciones con un código de estado y un mensaje de error",
            "type": "object",
            "properties": {
                "message": {
                    "description": "Mensaje de error\nexample \"Solicitud incorrecta\"",
                    "type": "string",
                    "example": "Solicitud incorrecta"
                },
                "status": {
                    "description": "Código de estado HTTP\nexample 400",
                    "type": "integer",
                    "example": 400
                }
            }
        }
    },
    "securityDefinitions": {
        "CookieAuth": {
            "type": "apiKey",
            "name": "Cookie",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "v1.0.0",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "GoGallery",
	Description:      "API para la gestión de subida de fotos, con una autentificación",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
