package convert

import (
	"errors"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/disintegration/imaging"
	pigo "github.com/esimov/pigo/core"
	facedetect "github.com/korylprince/go-face-detect"
	"golang.org/x/exp/slog"
)

func convertPortrait(c *config, inpath, outpath string) error {
	var (
		img *image.NRGBA
		err error
	)
	if c.useEXIF {
		img, err = facedetect.DecodeFileWithEXIF(inpath)
	} else {
		img, err = pigo.GetImage(inpath)
	}
	if err != nil {
		if img == nil {
			return fmt.Errorf("could not read image: %w", err)
		}
		c.logger.Debug("could not parse EXIF data", "input_path", inpath, "error", err)
	}

	portrait, err := c.detector.Portrait(img, c.portraitConfig)
	if err != nil {
		return err
	}

	if err = imaging.Save(portrait, outpath); err != nil {
		return fmt.Errorf("could not write portrait: %w", err)
	}

	return nil
}

type config struct {
	detector       *facedetect.Detector
	workers        int
	overwrite      bool
	useEXIF        bool
	logger         *slog.Logger
	portraitConfig *facedetect.PortraitConfig
}

type ConvertOption func(*config)

// WithWorkers configures the number concurrent workers.
// The default is runtime.NumCPU()
func WithWorkers(workers int) ConvertOption {
	return func(c *config) {
		c.workers = workers
	}
}

// WithOverwrite configures the converter to overwrite existing images.
// The default is false
func WithOverwrite(overwrite bool) ConvertOption {
	return func(c *config) {
		c.overwrite = overwrite
	}
}

// WithPortraitConfig configures the PortraitConfig for converting portraits.
// The default is facedetect.DefaultPortraitConfig
func WithPortraitConfig(pc *facedetect.PortraitConfig) ConvertOption {
	return func(c *config) {
		c.portraitConfig = pc
	}
}

// WithEXIF configures the converter to automatically rotate images based on EXIF orientation data.
// The default is true
func WithEXIF(useEXIF bool) ConvertOption {
	return func(c *config) {
		c.useEXIF = useEXIF
	}
}

// WithLogger configures the converter's logger.
// The default is a logger that logs error messages to stderr
func WithLogger(logger *slog.Logger) ConvertOption {
	return func(c *config) {
		c.logger = logger
	}
}

func worker(wg *sync.WaitGroup, c *config, outdir string, in chan string) {
	defer wg.Done()
	for inpath := range in {
		outpath := filepath.Join(outdir, filepath.Base(inpath))

		if _, err := os.Stat(outpath); !errors.Is(err, os.ErrNotExist) && !c.overwrite {
			c.logger.Debug("not overwriting existing file", "input_path", inpath, "output_path", outpath)
			continue
		} else if !errors.Is(err, os.ErrNotExist) && c.overwrite {
			c.logger.Debug("overwriting file", "input_path", inpath, "output_path", outpath)
		}
		if err := convertPortrait(c, inpath, outpath); err != nil {
			c.logger.Error("conversion failed", "input_path", inpath, "output_path", outpath, "error", err)
		} else {
			c.logger.Info("portrait converted", "input_path", inpath, "output_path", outpath)
		}
	}
}

// ConvertPortraits concurrently converts the images at paths given in infiles to portraits and outputs the results to outpath.
// It's recommended to use the embedded cascade.Detector. Check ConvertOption for configurable options
func ConvertPortraits(detector *facedetect.Detector, infiles []string, outdir string, opts ...ConvertOption) {
	c := &config{
		detector:       detector,
		workers:        runtime.NumCPU(),
		useEXIF:        true,
		logger:         slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})),
		portraitConfig: facedetect.DefaultPortraitConfig,
	}
	for _, opt := range opts {
		opt(c)
	}

	if err := os.MkdirAll(outdir, 0755); err != nil {
		c.logger.Error("could not create output directory", "path", outdir, "error", err)
		return
	}

	if c.workers > len(infiles) {
		c.workers = len(infiles)
	}

	in := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(c.workers)
	for i := 0; i < c.workers; i++ {
		go worker(wg, c, outdir, in)
	}

	for _, path := range infiles {
		in <- path
	}
	close(in)

	wg.Wait()
}
