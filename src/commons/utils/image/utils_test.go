package utilsImage

import (
	"go-gallery/src/commons/constants"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestResizeRealImages(t *testing.T) {
	// Directorio donde tienes las imágenes de ejemplo
	inputDir := "../../../test/resources/images"

	files, err := os.ReadDir(inputDir)
	if err != nil {
		t.Fatalf("Error leyendo el directorio de imágenes: %v", err)
	}

	// Definir las extensiones válidas para imágenes
	validExtensions := []string{constants.JPG_EXTENSION, constants.JPEG_EXTENSION, constants.WEBP_EXTENSION, constants.PNG_EXTENSION}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		ext := strings.ToLower(filepath.Ext(name))

		// Verificar si el archivo tiene una extensión válida
		if !contains(validExtensions, ext) {
			continue
		}

		inputPath := filepath.Join(inputDir, name)
		data, err := os.ReadFile(inputPath)
		if err != nil {
			t.Errorf("No se pudo leer la imagen %s: %v", name, err)
			continue
		}

		// Redimensionar la imagen
		webpBytes, err := ResizeImage(data, constants.THUMBNAIL_WIDTH, constants.THUMBNAIL_HEIGHT)
		if err != nil {
			t.Errorf("Error redimensionando %s: %v", name, err)
			continue
		}

		if len(webpBytes) == 0 {
			t.Errorf("La imagen redimensionada de %s está vacía", name)
		}

		if !strings.HasPrefix(string(webpBytes), "RIFF") {
			t.Errorf("La imagen redimensionada de %s no es un archivo WebP válido", name)
		}

		t.Logf("Imagen procesada y redimensionada: %s", inputPath)
	}
}

func TestResizeImagesError(t *testing.T) {
	data := []byte("testdata")
	_, err := ResizeImage(data, constants.THUMBNAIL_WIDTH, constants.THUMBNAIL_HEIGHT)
	if err == nil {
		t.Errorf("Se esperaba un error al intentar redimensionar la imagen")
		t.FailNow()
	}
}

func TestEncondeImageToBase64(t *testing.T) {
	data := []byte("testdata")
	b64 := EncondeImageToBase64(data)

	if b64 != "dGVzdGRhdGE=" {
		t.Errorf("El resultado al codificar la imagen a base64 es erróneo: got %s", b64)
	}
}

func TestHumanizeBytes(t *testing.T) {
	data := HumanizeBytes(1024)
	if data != "1.0 kB" {
		t.Errorf("El resultado al obtener el tamaño no es correcto got %s", data)
	}
}

// contains verifica si una cadena está en una lista de cadenas
func contains(list []string, str string) bool {
	return slices.Contains(list, str)
}
