package facedetect

import (
	"image"
	"image/color"
	"math"

	"github.com/disintegration/imaging"
)

// radToDegree converts radians to degrees
func radToDegree(rad float64) float64 {
	return (rad * 180) / math.Pi
}

// Rotate rotates the image so that the line going through the pupils is parallel with the top edge of the image
func Rotate(img image.Image, face *Face) *image.NRGBA {
	angle := radToDegree(math.Atan2(
		-float64(face.LeftEye.Row-face.RightEye.Row),
		-float64(face.LeftEye.Col-face.RightEye.Col),
	))

	return imaging.Rotate(img, angle, color.NRGBA{})
}

// checkCorners returns true if all four corners of the rectangle specified by center x, y and width and height
// are not transparent
func checkCorners(img image.Image, x, y, width, height int) bool {
	if _, _, _, a := img.At(x-(width/2), y-(height/2)).RGBA(); a == 0 {
		return false
	}
	if _, _, _, a := img.At(x-(width/2), y+(width/2)).RGBA(); a == 0 {
		return false
	}
	if _, _, _, a := img.At(x+(width/2), y-(height/2)).RGBA(); a == 0 {
		return false
	}
	if _, _, _, a := img.At(x+(width/2), y+(width/2)).RGBA(); a == 0 {
		return false
	}

	return true
}

func crop(img image.Image, x, y, width, height int) *image.NRGBA {
	rect := image.Rect(x-(width/2), y-(height/2), x+(width/2), y+(height/2))
	return imaging.Crop(img, rect)
}

// Crop crops the image to the largest bounding box that doesn't contain transparent pixels.
// aspectRatio is the ratio width / height.
// maxWidthRatio is the maximum width of the cropped image / the width of the detected face
func Crop(img image.Image, face *Face, aspectRatio, maxWidthRatio float64) *image.NRGBA {
	aspect := float64(1.0 / aspectRatio)
	minWidth := face.Bounds.Scale
	maxWidth := int(float64(minWidth) * maxWidthRatio)
	maxHeight := int(float64(maxWidth) * aspect)
	x := (face.LeftEye.Col + face.RightEye.Col) / 2
	y := (face.LeftEye.Row+face.RightEye.Row)/2 + int(float64(maxHeight)*0.1)

	if checkCorners(img, x, y, maxWidth, maxHeight) {
		return crop(img, x, y, maxWidth, maxHeight)
	}
	for {
		width := (maxWidth + minWidth) / 2
		height := int(float64(width) * aspect)

		if width == minWidth {
			return crop(img, x, y, width, height)
		}

		if checkCorners(img, x, y, width, height) {
			minWidth = width
		} else {
			maxWidth = width
		}
	}
}

// Brighten brightens the image for better detail in the face
func Brighten(img image.Image, brightness, contrast, gamma float64) *image.NRGBA {
	i := imaging.AdjustBrightness(img, brightness)
	i = imaging.AdjustContrast(i, contrast)
	i = imaging.AdjustGamma(i, gamma)
	return i
}
