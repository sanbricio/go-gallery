package utilsImage

import (
	"bytes"
	"encoding/base64"
	"image"

	_ "image/jpeg"  
	_ "image/png"      
	_ "golang.org/x/image/webp" 

	"github.com/HugoSmits86/nativewebp"
	"github.com/dustin/go-humanize"
	"golang.org/x/image/draw"
)

// ResizeImage redimensiona la imagen y la convierte a WebP.
func ResizeImage(input []byte, width, height int) ([]byte, error) {
	// Decodificar la imagen de entrada
	img, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	// Crear una nueva imagen en blanco para el redimensionamiento
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	// Redimensionar la imagen
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	var buf bytes.Buffer
	err = nativewebp.Encode(&buf, dst, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func EncondeImageToBase64(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

func HumanizeBytes(size uint64) string {
	return humanize.Bytes(size)
}
