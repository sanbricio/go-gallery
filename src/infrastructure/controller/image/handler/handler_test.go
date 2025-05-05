package imageHandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-gallery/src/commons/constants"
	"go-gallery/src/commons/exception"
	imageDTO "go-gallery/src/infrastructure/dto/image"
	"go-gallery/src/infrastructure/logger"
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
	{"Valid JPG", "../../../../test/resources/images/landscape.jpg", "landscape", constants.JPG_EXTENSION, false},
	{"Valid WEBP", "../../../../test/resources/images/landscape.webp", "landscape", constants.WEBP_EXTENSION, false},
	{"Valid PNG", "../../../../test/resources/images/landscape.png", "landscape", constants.PNG_EXTENSION, false},
	{"Valid JPEG", "../../../../test/resources/images/landscape.jpeg", "landscape", constants.JPEG_EXTENSION, false},
	{"Invalid TXT", "../../../../test/resources/images/landscape.txt", "landscape", "", true},
}

func beforeAll() {
	logger.Init(logger.NewConsoleLogger())
}

func TestProcessImageFile(t *testing.T) {
	beforeAll()
	app := loadFiberApp()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := createRequest(t, tt.imagePath)

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Error in test request: %v", err)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error reading the response: %v", err)
			}

			if tt.expectError {
				evaluateWrongImage(t, resp, body)
				return
			}

			evaluateImage(t, body, tt.expectedName, tt.expectedExt)

		})
	}
}

func TestProcessImageFileFailToOpen(t *testing.T) {
	beforeAll()
	fileHeader := &multipart.FileHeader{
		Filename: "nonexistent.jpg",
	}

	result, apiErr := ProcessImageFile(fileHeader, "testOwner")

	if result != nil {
		t.Errorf("Expected result to be nil, but got: %+v", result)
	}
	if apiErr == nil || apiErr.Message != "Error opening the image file" {
		t.Errorf("Expected error when opening the file, but got: %v", apiErr)
	}
}

func TestEncodeToRawBytes(t *testing.T) {
	beforeAll()
	fileHeader := &multipart.FileHeader{
		Filename: "nonexistent.jpg",
	}
	_, apiErr := encodeToRawBytes(fileHeader)

	if apiErr == nil || apiErr.Message != "Error opening the image file" {
		t.Errorf("Expected error opening the file, but got: %v", apiErr)
	}
}

func TestReadAllFile(t *testing.T) {
	beforeAll()
	reader := &failingFile{}

	_, apiErr := readAllFile(reader)
	if apiErr == nil || apiErr.Message != "Error reading the image file" {
		t.Errorf("Expected error reading the file, but got: %v", apiErr)
	}
}

func loadFiberApp() *fiber.App {
	// Initialize Fiber
	app := fiber.New()

	// Define test route
	app.Post("/test", func(c *fiber.Ctx) error {
		// Get the file from the form
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Could not read the file"})
		}

		// Process the image with the function we are testing
		result, apiErr := ProcessImageFile(file, "testOwner")
		if apiErr != nil {
			return c.Status(apiErr.Status).JSON(apiErr)
		}

		return c.JSON(result)
	})
	return app
}

func evaluateImage(t *testing.T, body []byte, expectedName, expectedExt string) {
	// Verify expected response
	var result imageDTO.ImageUploadRequestDTO
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Error parsing the JSON response: %v", err)
	}

	// Validate correct data
	if result.Name != expectedName {
		t.Errorf("Expected name '%s', but got '%s'", expectedName, result.Name)
	}
	if result.Extension != expectedExt {
		t.Errorf("Expected extension '%s', but got '%s'", expectedExt, result.Extension)
	}
	if result.Owner != "testOwner" {
		t.Errorf("Expected Owner 'testOwner', but got '%s'", result.Owner)
	}
	if result.Size == "" {
		t.Errorf("The 'Size' field is empty")
	}
	if len(result.RawContentFile) == 0 {
		t.Errorf("The 'RawContentFile' field is empty")
	}
}

func evaluateWrongImage(t *testing.T, resp *http.Response, body []byte) {
	// Verify that the expected error is returned
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, but got %d", resp.StatusCode)
	}

	var apiErr exception.ApiException
	if err := json.Unmarshal(body, &apiErr); err != nil {
		t.Fatalf("Error parsing the error JSON response: %v", err)
	}

	expectedMsg := "Unsupported file format. Only jpg, jpeg, png, and webp images are accepted."
	if apiErr.Message != expectedMsg {
		t.Errorf("Incorrect error message: %s", apiErr.Message)
	}
}

func createRequest(t *testing.T, imagePath string) *http.Request {
	file, err := os.Open(imagePath)
	if err != nil {
		t.Fatalf("Failed to open test image '%s': %v", imagePath, err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		t.Fatalf("Error creating form file part: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatalf("Error copying file into form: %v", err)
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/test", &requestBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

// Mock to simulate error when reading file

type failingFile struct{}

func (f *failingFile) Read(p []byte) (n int, err error) {
	return 0, errors.New("forced read error")
}

func (f *failingFile) Close() error {
	return nil
}

// Additional methods to implement multipart.File but not used in this test
func (f *failingFile) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (f *failingFile) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, io.EOF
}
