package image

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/chai2010/webp"
	"github.com/thediveo/enumflag/v2"
)

type encoder func(io.Writer, image.Image) error

type imgEncoding struct {
	encoder   encoder
	extension string
	format    ImgFormat
}

func (e imgEncoding) Extension() string {
	return e.extension
}

func (e imgEncoding) Encode(w io.Writer, i image.Image) error {
	return e.encoder(w, i)
}

var (
	jpegEncoding imgEncoding = imgEncoding{format: JPEG, extension: "jpg", encoder: func(w io.Writer, i image.Image) error { return jpeg.Encode(w, i, nil) }}
	pngEncoding  imgEncoding = imgEncoding{format: PNG, extension: "png", encoder: png.Encode}
	webpEncoding imgEncoding = imgEncoding{format: WEBP, extension: "webp", encoder: func(w io.Writer, i image.Image) error { return webp.Encode(w, i, nil) }}
)

func EncodingFromFormatName(format string) (imgEncoding, error) {
	foundFormat := UNSUPPORTED
	for key, val := range formatNames {
		if val[0] == format {
			foundFormat = key
			break
		}
	}
	return EncodingFromFormat(foundFormat)
}

func EncodingFromFormat(format ImgFormat) (imgEncoding, error) {
	switch format {
	case jpegEncoding.format:
		return jpegEncoding, nil
	case pngEncoding.format:
		return pngEncoding, nil
	case webpEncoding.format:
		return webpEncoding, nil
	default:
		return imgEncoding{}, fmt.Errorf("unsupported image format '%v'", format)
	}
}

type ImgFormat enumflag.Flag

const (
	UNSUPPORTED ImgFormat = iota
	JPEG
	PNG
	WEBP
)

type formatIdMapping map[ImgFormat][]string

var formatNames = formatIdMapping{
	JPEG: {"jpeg"},
	PNG:  {"png"},
	WEBP: {"webp"},
}

func GetFormatNames() map[ImgFormat][]string {
	newFormats := make(formatIdMapping, 3)
	for k, v := range formatNames {
		newFormats[k] = v
	}
	return newFormats
}
