package image

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/chai2010/webp"
)

type encoder func(io.Writer, image.Image) error

type imgEncoding struct {
	encoder   encoder
	format    string
	extension string
}

func (e imgEncoding) Extension() string {
	return e.extension
}

func (e imgEncoding) Encode(w io.Writer, i image.Image) error {
	return e.encoder(w, i)
}

var (
	jpegEncoding imgEncoding = imgEncoding{format: "jpeg", extension: "jpg", encoder: func(w io.Writer, i image.Image) error { return jpeg.Encode(w, i, nil) }}
	pngEncoding  imgEncoding = imgEncoding{format: "png", extension: "png", encoder: png.Encode}
	webpEncoding imgEncoding = imgEncoding{format: "webp", extension: "webp", encoder: func(w io.Writer, i image.Image) error { return webp.Encode(w, i, nil) }}
)

func EncodingFromFormat(format string) (imgEncoding, error) {
	switch format {
	case jpegEncoding.format:
		return jpegEncoding, nil
	case pngEncoding.format:
		return pngEncoding, nil
	case webpEncoding.format:
		return webpEncoding, nil
	default:
		return imgEncoding{}, fmt.Errorf("unsupported image format '%s'", format)
	}
}
