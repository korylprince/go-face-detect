package facedetect

import (
	"errors"
	"image"

	pigo "github.com/esimov/pigo/core"
)

var (
	ErrFaceUndetected   = errors.New("face undetected")
	ErrPupilsUndetected = errors.New("pupils undetected")
)

// DetectParams are the parameters given to the face detector.
// MinSizeFactor is the smallest size area searched for as a percentage of the largest dimension of the image.
// MaxSizeFactor is the largest size area searched for as a percentage of the largest dimension of the image.
// ShiftFactor determines to what percentage to move the detection window over its size.
// ScaleFactor defines in percentage the resize value of the detection window when moving to a higher scale.
// IoUThreshold is the threshold to consider multiple face regions the same face
type DetectParams struct {
	MinSizeFactor float64
	MaxSizeFactor float64
	ShiftFactor   float64
	ScaleFactor   float64
	IoUThreshold  float64
}

var FastDetectParams = &DetectParams{
	MinSizeFactor: 0.2,
	MaxSizeFactor: 0.8,
	ShiftFactor:   0.15,
	ScaleFactor:   1.15,
	IoUThreshold:  0.15,
}

var SlowDetectParams = &DetectParams{
	MinSizeFactor: 0.1,
	MaxSizeFactor: 0.9,
	ShiftFactor:   0.05,
	ScaleFactor:   1.03,
	IoUThreshold:  0,
}

type Detector struct {
	FaceCascade  *pigo.Pigo
	PupilCascade *pigo.PuplocCascade
}

type Face struct {
	Bounds   pigo.Detection
	LeftEye  *pigo.Puploc
	RightEye *pigo.Puploc
}

func (d *Detector) DetectFaces(img pigo.ImageParams, params *DetectParams, angle float64) []pigo.Detection {
	maxSize := img.Rows
	if img.Cols > img.Rows {
		maxSize = img.Cols
	}
	p := pigo.CascadeParams{
		MinSize:     int(params.MinSizeFactor * float64(maxSize)),
		MaxSize:     int(params.MaxSizeFactor * float64(maxSize)),
		ShiftFactor: params.ShiftFactor,
		ScaleFactor: params.ScaleFactor,
		ImageParams: img,
	}

	// find all faces
	faces := d.FaceCascade.RunCascade(p, angle)
	// filter duplicate faces
	return d.FaceCascade.ClusterDetections(faces, params.IoUThreshold)
}

// ChooseBestFace returns the face with the highest quality (Q)
func ChooseBestFace(faces []pigo.Detection) pigo.Detection {
	if len(faces) == 1 {
		return faces[0]
	}

	best := faces[0]
	for _, face := range faces[1:] {
		if face.Q > best.Q {
			best = face
		}
	}
	return best
}

func (d *Detector) DetectPupils(img pigo.ImageParams, face pigo.Detection, angle float64) (leftEye, rightEye *pigo.Puploc) {
	// search in general area of left pupil
	puploc := &pigo.Puploc{
		Row:      face.Row - int(0.085*float32(face.Scale)),
		Col:      face.Col - int(0.185*float32(face.Scale)),
		Scale:    float32(face.Scale) * 0.4,
		Perturbs: 50,
	}
	leftEye = d.PupilCascade.RunDetector(*puploc, img, angle, false)

	// search in general area of right pupil
	puploc = &pigo.Puploc{
		Row:      face.Row - int(0.085*float32(face.Scale)),
		Col:      face.Col + int(0.185*float32(face.Scale)),
		Scale:    float32(face.Scale) * 0.4,
		Perturbs: 50,
	}
	rightEye = d.PupilCascade.RunDetector(*puploc, img, angle, false)

	return leftEye, rightEye
}

// DetectFace detects a single face and pupils in an image, returning the detected areas.
// DetectFace attempts detection using FastDetectParams, falling back to SlowDetectParams if a face isn't detected
func (d *Detector) DetectFace(img *image.NRGBA, angle float64) (*Face, error) {
	x, y := img.Bounds().Max.X, img.Bounds().Max.Y
	params := pigo.ImageParams{
		Pixels: pigo.RgbToGrayscale(img),
		Cols:   x,
		Rows:   y,
		Dim:    x,
	}

	// try to detect faces with faster detection first and fallback to slower detection if it fails
	faces := d.DetectFaces(params, FastDetectParams, angle)
	if len(faces) == 0 {
		faces = d.DetectFaces(params, SlowDetectParams, angle)
		if len(faces) == 0 {
			return nil, ErrFaceUndetected
		}
	}

	face := &Face{Bounds: ChooseBestFace(faces)}

	face.LeftEye, face.RightEye = d.DetectPupils(params, face.Bounds, angle)
	if face.LeftEye.Row <= 0 || face.LeftEye.Col <= 0 || face.RightEye.Row <= 0 || face.RightEye.Col <= 0 {
		return face, ErrPupilsUndetected
	}

	return face, nil
}
