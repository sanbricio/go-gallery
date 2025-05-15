package utilsImage

import (
	"go-gallery/src/commons/constants"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResizeRealImages(t *testing.T) {
	// Directorio donde tienes las imágenes de ejemplo
	inputDir := "../../../test/resources/images"

	files, err := os.ReadDir(inputDir)
	require.NoError(t, err, "Error leyendo el directorio de imágenes")

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
		assert.NoError(t, err, "No se pudo leer la imagen %s", name)

		// Redimensionar la imagen
		webpBytes, err := ResizeImage(data, constants.THUMBNAIL_WIDTH, constants.THUMBNAIL_HEIGHT)
		assert.NoError(t, err, "Error redimensionando %s", name)

		assert.NotEmpty(t, webpBytes, "La imagen redimensionada de %s está vacía", name)
		assert.True(t, strings.HasPrefix(string(webpBytes), "RIFF"), "La imagen redimensionada de %s no es un archivo WebP válido", name)

		t.Logf("Imagen procesada y redimensionada: %s", inputPath)
	}
}

func TestResizeImagesError(t *testing.T) {
	data := []byte("testdata")
	_, err := ResizeImage(data, constants.THUMBNAIL_WIDTH, constants.THUMBNAIL_HEIGHT)
	assert.Error(t, err, "Se esperaba un error al intentar redimensionar la imagen")
}

func TestEncondeImageToBase64(t *testing.T) {
	data := []byte("testdata")
	b64 := EncondeImageToBase64(data)

	assert.Equal(t, "dGVzdGRhdGE=", b64, "El resultado al codificar la imagen a base64 es erróneo")
}

func TestHumanizeBytes(t *testing.T) {
	data := HumanizeBytes(1024)
	assert.Equal(t, "1.0 kB", data, "El resultado al obtener el tamaño no es correcto")
}

// contains verifica si una cadena está en una lista de cadenas
func contains(list []string, str string) bool {
	return slices.Contains(list, str)
}
