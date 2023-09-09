package facedetect

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging"
	pigo "github.com/esimov/pigo/core"
)

type PortraitConfig struct {
	AspectRatio   float64
	MaxWidthRatio float64
	Brightness    float64
	Contrast      float64
	Gamma         float64
}

var DefaultPortraitConfig = &PortraitConfig{
	AspectRatio:   3.0 / 4.0,
	MaxWidthRatio: 1.5,
	Brightness:    0,
	Contrast:      5,
	Gamma:         1.4,
}

// Portrait detects a single face in an image, rotates, crops, and brightens it, and returns the result.
// If config is nil, DefaultPortraitConfig is used
func (d *Detector) Portrait(img *image.NRGBA, config *PortraitConfig) (*image.NRGBA, error) {
	if config == nil {
		config = DefaultPortraitConfig
	}

	// detect face
	face, err := d.DetectFace(img, 0)
	if err != nil {
		return nil, fmt.Errorf("could not detect face: %w", err)
	}

	// rotate based on pupils
	rotated := Rotate(img, face)

	// detect rotated face
	face, err = d.DetectFace(rotated, 0)
	if err != nil {
		return nil, fmt.Errorf("could not detect rotated face: %w", err)
	}

	cropped := Crop(rotated, face, config.AspectRatio, config.MaxWidthRatio)
	brightened := Brighten(cropped, config.Brightness, config.Contrast, config.Gamma)

	return brightened, nil
}

// PortraitFile detects a single face in the image at inpath, rotates, crops, and brightens it, and writes the result to outpath.
// If config is nil, DefaultPortraitConfig is used
func (d *Detector) PortraitFile(inpath, outpath string, config *PortraitConfig) error {
	img, err := pigo.GetImage(inpath)
	if err != nil {
		return fmt.Errorf("could not open image %s: %w", inpath, err)
	}

	img, err = d.Portrait(img, config)
	if err != nil {
		return fmt.Errorf("could not convert image: %w", err)
	}

	if err = imaging.Save(img, outpath); err != nil {
		return fmt.Errorf("could not write portrait: %w", err)
	}

	return nil
}

// PortraitFileWithEXIF reads EXIF data from the image at inpath, rotating it if necessary,
// detects a single face in the image, rotates, crops, and brightens it, and writes the result to outpath.
// If config is nil, DefaultPortraitConfig is used
func (d *Detector) PortraitFileWithEXIF(inpath, outpath string, config *PortraitConfig) error {
	img, err := DecodeFileWithEXIF(inpath)
	if err != nil {
		if img == nil {
			return fmt.Errorf("could not open image %s: %w", inpath, err)
		}
	}

	img, err = d.Portrait(img, config)
	if err != nil {
		return fmt.Errorf("could not convert image: %w", err)
	}

	if err = imaging.Save(img, outpath); err != nil {
		return fmt.Errorf("could not write portrait: %w", err)
	}

	return nil
}
