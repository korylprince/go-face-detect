package cascade

import (
	_ "embed"
	"fmt"

	pigo "github.com/esimov/pigo/core"
	facedetect "github.com/korylprince/go-face-detect"
)

func unpackFaceFinder(cascade []byte) *pigo.Pigo {
	p := pigo.NewPigo()
	classifier, err := p.Unpack(cascade)
	if err != nil {
		panic(fmt.Errorf("could not unpack facefinder cascade: %w", err))
	}
	return classifier
}

//go:embed facefinder
var faceFinderCascade []byte
var FaceCascade = unpackFaceFinder(faceFinderCascade)

func unpackPupilLocator(cascade []byte) *pigo.PuplocCascade {
	p := pigo.NewPuplocCascade()
	classifier, err := p.UnpackCascade(cascade)
	if err != nil {
		panic(fmt.Errorf("could not unpack puploc cascade: %w", err))
	}
	return classifier
}

//go:embed puploc
var pupilLocatorCascade []byte
var PupilCascade = unpackPupilLocator(pupilLocatorCascade)

var Detector = &facedetect.Detector{
	FaceCascade:  FaceCascade,
	PupilCascade: PupilCascade,
}
