package facedetect

import (
	"errors"
	"fmt"
	"image"
	"io"
	"os"

	"github.com/disintegration/imaging"
	pigo "github.com/esimov/pigo/core"
	"github.com/mholt/goexif2/exif"
	"github.com/mholt/goexif2/tiff"
)

var (
	ErrNoOrientationTag      = errors.New("no orientation tag")
	ErrInvalidOrientationTag = errors.New("invalid orientation tag")
)

func decodeEXIF(r io.Reader) (int, error) {
	exifData, err := exif.Decode(r)
	if err != nil {
		return 0, fmt.Errorf("could not parse exif data: %w", err)
	}

	tag, err := exifData.Get(exif.Orientation)
	if err != nil {
		tagErr := new(exif.TagNotPresentError)
		if errors.As(err, tagErr) {
			return 0, ErrNoOrientationTag
		}
		return 0, fmt.Errorf("could not get orientation flag: %w", err)
	}

	if tag.Format() != tiff.IntVal || tag.Count < 1 {
		return 0, ErrInvalidOrientationTag
	}

	return tag.Int(0)
}

func rotateEXIF(img *image.NRGBA, orientation int) *image.NRGBA {
	switch orientation {
	case 2:
		return imaging.FlipH(img)
	case 3:
		return imaging.Rotate180(img)
	case 4:
		return imaging.Rotate180(imaging.FlipH(img))
	case 5:
		return imaging.Rotate270(imaging.FlipV(img))
	case 6:
		return imaging.Rotate270(img)
	case 7:
		return imaging.Rotate90(imaging.FlipV(img))
	case 8:
		return imaging.Rotate90(img)
	}

	return img
}

// DecodeWithEXIF decodes an image from r, rotating it using EXIF data if it exists
func DecodeWithEXIF(r io.ReadSeeker) (*image.NRGBA, error) {
	img, err := pigo.DecodeImage(r)
	if err != nil {
		return nil, fmt.Errorf("could not decode image: %w", err)
	}

	if _, err = r.Seek(0, io.SeekStart); err != nil {
		return img, fmt.Errorf("could not seek to start of reader: %w", err)
	}

	orientation, err := decodeEXIF(r)
	if err != nil {
		return img, fmt.Errorf("could not get exif orientation: %w", err)
	}

	return rotateEXIF(img, orientation), nil
}

// DecodeFileWithEXIF decodes the image at path, rotating it using EXIF data if it exists
func DecodeFileWithEXIF(path string) (*image.NRGBA, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %w", path, err)
	}
	defer f.Close()

	return DecodeWithEXIF(f)
}
