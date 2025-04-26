package utils

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"

	"golang.org/x/image/draw"
)

func ResizeImage(input []byte, width, height int) ([]byte, string, error) {
	img, format, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, "", err
	}

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	var buf bytes.Buffer
	switch format {
	case "jpeg":
		err = jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 80})
	case "png":
		err = png.Encode(&buf, dst)
	default:
		return nil, "", err
	}

	if err != nil {
		return nil, "", err
	}

	return buf.Bytes(), format, nil
}
