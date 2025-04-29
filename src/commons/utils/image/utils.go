package utilsImage

import (
	"bytes"
	"encoding/base64"
	"go-gallery/src/commons/constants"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/dustin/go-humanize"
	"golang.org/x/image/draw"
)

func ResizeImage(input []byte, width, height int) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	var buf bytes.Buffer
	//TODO AÑADIR los formatos webp y png creo que faltan mirar el handler
	switch format {
	case constants.JPEG:
		err = jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 80})
	case constants.PNG:
		err = png.Encode(&buf, dst)
	default:
		return nil, err
	}

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
