[![Go Reference](https://pkg.go.dev/badge/github.com/korylprince/go-face-detect.svg)](https://pkg.go.dev/github.com/korylprince/go-face-detect)

# About

go-face-detect is a Go library, cli tool, and Web Assembly module to detect faces in images (using [pigo](https://github.com/esimov/pigo)), and transform them for use as a portrait (e.g. for a name badge).

# Library

The [library](https://pkg.go.dev/github.com/korylprince/go-face-detect) contains high-level functions to configurably transform images into portraits, as well as low-level general face detection and image transformation functions.

# CLI tool

The included [face-detect](https://github.com/korylprince/go-face-detect/tree/master/cmd/face-detect/) utility can quickly process multiple images in parallel, applying the following transformations:

* Rotate the image so face is level
* Crop the image so the face is well framed
* Brighten the image so face detail is easier to see

```
Usage: face-detect [flags] -out <output directory> <input file>...
  -aspect-ratio float
    	the width / height aspect ratio for the converted portraits (default 0.75)
  -brightness float
    	the percentage to adjust the converted portrait brightness (-100 to 100)
  -contrast float
    	the percentage to adjust the converted portrait contrast (-100 to 100) (default 5)
  -gamma float
    	the amount to adjust the converted portrait gamma (1.0 returns the gamma as-is) (default 1.4)
  -level string
    	logging level parsable by slog.UnmarshalText (default "INFO")
  -max-width-ratio float
    	the max portrait width / detected face width ratio (default 1.5)
  -out string
    	the directory where converted portraits will be written
  -overwrite
    	overwrite existing files
  -use-exif
    	automatically rotate photos based on EXIF orientation (default true)
  -workers int
    	number of concurrent workers to use (default 16)

face-detect detects a single face in an image, automatically rotates, crops, brightens the image and writes it to a new file.
If multiple input images are given, they'll be processed in parallel.
```

# Web Assembly Live Demo

go-face-detect also includes a [wasm](https://github.com/korylprince/go-face-detect/tree/master/wasm/) module that can be compiled to a standalone wasm file (including embedded cascade files).

Visit [https://korylprince.github.io/go-face-detect/](https://korylprince.github.io/go-face-detect/) for a live demo using wasm to convert portraits in the browser.

![Screenshot](https://raw.githubusercontent.com/korylprince/go-face-detect/master/screenshot.png)

# Dependencies

* [github.com/esimov/pigo](https://github.com/esimov/pigo) for pure-Go face detection
* [github.com/mholt/goexif2](https://github.com/mholt/goexif2) for EXIF parsing
* [github.com/disintegration/imaging](https://github.com/disintegration/imaging) for image manipulation

# License

`facefinder` and `puploc` cascade files are from the [pigo](https://github.com/esimov/pigo). These were originally created by [Nenad Markuš](https://github.com/nenadmarkus).

`wasm_exec.js` is included in Go's distribution.

All other code licensed by [LICENSE](https://github.com/korylprince/go-face-detect/blob/master/LICENSE).
