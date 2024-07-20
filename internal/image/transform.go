package image

import (
	"fmt"
	"image"
	"path/filepath"
	"runtime"

	"github.com/eslam-allam/image-tweaker/internal/fs"
	"github.com/nfnt/resize"
)

func ResizeIfBigger(img image.Image, width, height uint) image.Image {
	bounds := img.Bounds()
	currentWidth, currentHeight := bounds.Dx(), bounds.Dy()
	if int(width) > currentWidth {
		width = uint(currentWidth)
	}
	if int(height) > currentHeight {
		height = uint(currentHeight)
	}
	if currentHeight == int(height) && currentWidth == int(width) {
		return img
	}
	return resize.Resize(width, height, img, resize.Lanczos3)
}

func Transform(srcPath, dstPath string, resize bool, targetWidth, targetHeight uint, targetFormat ImgFormat, threads uint) error {
	var err error
	if !filepath.IsAbs(srcPath) {
		srcPath, err = filepath.Abs(srcPath)
		if err != nil {
			return err
		}
	}

	if !filepath.IsAbs(dstPath) {
		dstPath, err = filepath.Abs(dstPath)
		if err != nil {
			return err
		}
	}

	srcExists, srcTyp, err := fs.Exists(srcPath)
	if err != nil {
		return err
	}
	if !srcExists {
		return fmt.Errorf("path '%s' does not exists", srcPath)
	}

	if srcTyp == fs.File {
		return transformFile(srcPath, dstPath, resize, targetWidth, targetHeight, targetFormat)
	} else {
		if threads == 0 {
			threads = uint(runtime.NumCPU())
		}

		return fs.ApplyDirectoryParallel(threads, srcPath, dstPath, func(fileSrc, fileDst string) error {
			return transformFile(fileSrc, fileDst, resize, targetWidth, targetHeight, targetFormat)
		})
	}
}

func transformFile(srcPath, dstPath string, resize bool, targetWidth uint, targetHeight uint, targetFormat ImgFormat) error {
	img, format, err := ReadImage(srcPath)
	if err != nil {
		return err
	}

	if resize {
		img = ResizeIfBigger(img, targetWidth, targetHeight)
	}

	enc := format
	if targetFormat != UNSUPPORTED {
		enc, err = EncodingFromFormat(targetFormat)
		if err != nil {
			return err
		}
	}

	exists, typ, err := fs.Exists(dstPath)
	if err != nil {
		return err
	}
	if exists && typ == fs.Directory {
		dstPath = filepath.Join(dstPath, filepath.Base(srcPath)+"."+enc.extension)
	} else if filepath.Ext(dstPath) == "" {
		dstPath = dstPath + "." + enc.extension
	}
	if exists, _, err = fs.Exists(dstPath); err != nil || exists {
		return fmt.Errorf("output file '%s' already exists", dstPath)
	}

	err = SaveImage(img, enc, dstPath)
	if err != nil {
		return err
	}
	return nil
}
