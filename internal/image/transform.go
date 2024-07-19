package image

import (
	"image"

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
