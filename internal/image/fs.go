package image

import (
	"errors"
	"fmt"
	"image"
	"os"
)

type pathtype int

const (
	File pathtype = iota
	Directory
)

func exists(path string) (bool, pathtype, error) {
	info, err := os.Stat(path)
	if err == nil {
		var t pathtype
		if info.IsDir() {
			t = Directory
		} else {
			t = File
		}
		return true, t, nil
	}
	if os.IsNotExist(err) {
		return false, 0, nil
	}
	return false, 0, err
}

func ReadImage(path string) (image.Image, imgEncoding, error) {
	exists, typ, err := exists(path)
	if err != nil {
		return nil, imgEncoding{}, err
	}
	if !exists {
		return nil, imgEncoding{}, errors.New("specified path does not exist")
	}
	if typ != File {
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
	exists, _, err := exists(path)
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
