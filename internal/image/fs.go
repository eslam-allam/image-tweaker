package image

import (
	"errors"
	"fmt"
	"image"
	"os"

	"github.com/eslam-allam/image-tweaker/internal/fs"
)

func ReadImage(path string) (image.Image, imgEncoding, error) {
	exists, typ, err := fs.Exists(path)
	if err != nil {
		return nil, imgEncoding{}, err
	}
	if !exists {
		return nil, imgEncoding{}, errors.New("specified path does not exist")
	}
	if typ != fs.File {
		return nil, imgEncoding{}, errors.New("path must refer to a file")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, imgEncoding{}, fmt.Errorf("failed to open image '%s' for reading: %w", path, err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, imgEncoding{}, fmt.Errorf("failed to decode image: '%s': %w", path, err)
	}
	supportedFormat, err := EncodingFromFormatName(format)
	if err != nil {
		return nil, imgEncoding{}, err
	}
	return img, supportedFormat, nil
}

func SaveImage(img image.Image, encoding imgEncoding, path string) error {
	exists, _, err := fs.Exists(path)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("failed to save image to '%s': path already exists", path)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = encoding.Encode(file, img)
	if err != nil {
		return err
	}
	return nil
}

func IsExtensionSupported(ext string) bool {
	_, err := getEncodingFromExtension(ext)
	return err == nil
}
