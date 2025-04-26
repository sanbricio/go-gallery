package imageHandler

import (
	"bytes"
	"encoding/json"
	"go-gallery/src/commons/exception"
	imageDTO "go-gallery/src/infrastructure/dto/image"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
)

var tests = []struct {
	name         string
	imagePath    string
	expectedName string
	expectedExt  string
	expectError  bool
}{
	{"Valid JPG", "../../../resources/images/landscape.jpg", "landscape", ".jpg", false},
	{"Valid WEBP", "../../../resources/images/landscape.webp", "landscape", ".webp", false},
	{"Valid PNG", "../../../resources/images/landscape.png", "landscape", ".png", false},
	{"Valid JPEG", "../../../resources/images/landscape.jpeg", "landscape", ".jpeg", false},
	{"Invalid TXT", "../../../resources/images/landscape.txt", "landscape", "", true},
}

func TestProcessImageFile(t *testing.T) {
	app := loadFiberApp()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := createRequest(t, tt.imagePath)

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Error en la solicitud de prueba: %v", err)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error al leer la respuesta: %v", err)
			}

			if tt.expectError {
				evaluateWrongImage(t, resp, body)
				return
			}

			evaluateImage(t, body, tt.expectedName, tt.expectedExt)

		})
	}
}

func loadFiberApp() *fiber.App {
	// Inicializar Fiber
	app := fiber.New()

	// Definir la ruta de prueba
	app.Post("/test", func(c *fiber.Ctx) error {
		// Obtener el archivo desde el formulario
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No se pudo leer el archivo"})
		}

		// Procesar la imagen con la función que estamos probando
		result, apiErr := ProcessImageFile(file, "testOwner")
		if apiErr != nil {
			return c.Status(apiErr.Status).JSON(apiErr)
		}

		return c.JSON(result)
	})
	return app
}

func evaluateImage(t *testing.T, body []byte, expectedName, expectedExt string) {
	// Verificar respuesta esperada
	var result imageDTO.ImageDTO
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Error al parsear la respuesta JSON: %v", err)
	}

	// Validar que los datos sean correctos
	if result.Name != expectedName {
		t.Errorf("Se esperaba el nombre '%s', pero se obtuvo '%s'", expectedName, result.Name)
	}
	if result.Extension != expectedExt {
		t.Errorf("Se esperaba la extensión '%s', pero se obtuvo '%s'", expectedExt, result.Extension)
	}
	if result.Owner != "testOwner" {
		t.Errorf("Se esperaba el Owner 'testOwner', pero se obtuvo '%s'", result.Owner)
	}
	if result.Size == "" {
		t.Errorf("El campo Size está vacío")
	}
	if result.ContentFile == "" {
		t.Errorf("El campo ContentFile está vacío")
	}

}

func evaluateWrongImage(t *testing.T, resp *http.Response, body []byte) {
	// Verificar que se devuelva un error esperado
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Se esperaba código 400, pero se obtuvo %d", resp.StatusCode)
	}

	var apiErr exception.ApiException
	if err := json.Unmarshal(body, &apiErr); err != nil {
		t.Fatalf("Error al parsear la respuesta de error JSON: %v", err)
	}

	expectedMsg := "Formato de archivo no soportado. Solo se aceptan imágenes jpg, jpeg, png y webp"
	if apiErr.Message != expectedMsg {
		t.Errorf("Mensaje de error incorrecto: %s", apiErr.Message)
	}

}

func createRequest(t *testing.T, imagePath string) *http.Request {
	file, err := os.Open(imagePath)
	if err != nil {
		t.Fatalf("No se pudo abrir la imagen de prueba '%s': %v", imagePath, err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		t.Fatalf("Error al crear la parte del archivo en el formulario: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatalf("Error al copiar el archivo en el formulario: %v", err)
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/test", &requestBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}
