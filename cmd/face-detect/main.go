package main

import (
	"flag"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"runtime"

	facedetect "github.com/korylprince/go-face-detect"
	"github.com/korylprince/go-face-detect/cascade"
	convert "github.com/korylprince/go-face-detect/converter"
	"golang.org/x/exp/slog"
)

var Usage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] -out <output directory> <input file>...\n", filepath.Base(os.Args[0]))
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\n%s detects a single face in an image, automatically rotates, crops, brightens the image and writes it to a new file.\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(flag.CommandLine.Output(), "If multiple input images are given, they'll be processed in parallel.\n")
}

func main() {
	flWorkers := flag.Int("workers", runtime.NumCPU(), "number of concurrent workers to use")
	flOverwrite := flag.Bool("overwrite", false, "overwrite existing files")
	flUseEXIF := flag.Bool("use-exif", true, "automatically rotate photos based on EXIF orientation")
	flLogLevel := flag.String("level", "INFO", "logging level parsable by slog.UnmarshalText")
	flOutPath := flag.String("out", "", "the directory where converted portraits will be written")
	flAspectRatio := flag.Float64("aspect-ratio", 3.0/4.0, "the width / height aspect ratio for the converted portraits")
	flMaxWidthRatio := flag.Float64("max-width-ratio", 1.5, "the max portrait width / detected face width ratio")
	flBrightness := flag.Float64("brightness", 0, "the percentage to adjust the converted portrait brightness (-100 to 100)")
	flContrast := flag.Float64("contrast", 5, "the percentage to adjust the converted portrait contrast (-100 to 100)")
	flGamma := flag.Float64("gamma", 1.4, "the amount to adjust the converted portrait gamma (1.0 returns the gamma as-is)")

	flag.Usage = Usage
	flag.Parse()
	infiles := flag.Args()

	if len(infiles) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if *flOutPath == "" {
		fmt.Println("-out must be given")
		flag.Usage()
		os.Exit(1)
	}

	portraitConfig := &facedetect.PortraitConfig{
		AspectRatio:   *flAspectRatio,
		MaxWidthRatio: *flMaxWidthRatio,
		Brightness:    *flBrightness,
		Contrast:      *flContrast,
		Gamma:         *flGamma,
	}

	level := new(slog.Level)
	if err := level.UnmarshalText([]byte(*flLogLevel)); err != nil {
		fmt.Printf("could not parse -level (%s): %v\n", *flLogLevel, err)
		flag.Usage()
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: *level}))

	opts := []convert.ConvertOption{
		convert.WithWorkers(*flWorkers),
		convert.WithOverwrite(*flOverwrite),
		convert.WithPortraitConfig(portraitConfig),
		convert.WithEXIF(*flUseEXIF),
		convert.WithLogger(logger),
	}

	convert.ConvertPortraits(cascade.Detector, infiles, *flOutPath, opts...)
}
