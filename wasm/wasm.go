package main

import (
	"bytes"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"syscall/js"

	"github.com/disintegration/imaging"
	facedetect "github.com/korylprince/go-face-detect"
	"github.com/korylprince/go-face-detect/cascade"
)

func convert(buf []byte) ([]byte, error) {
	img, err := facedetect.DecodeWithEXIF(bytes.NewReader(buf))
	if err != nil {
		if img == nil {
			return nil, fmt.Errorf("could not decode image: %w", err)
		}
	}

	img, err = cascade.Detector.Portrait(img, facedetect.DefaultPortraitConfig)
	if err != nil {
		return nil, fmt.Errorf("could not convert to portrait: %w", err)
	}

	out := new(bytes.Buffer)
	if err = imaging.Encode(out, img, imaging.PNG); err != nil {
		return nil, fmt.Errorf("could not encode to png: %w", err)
	}

	return out.Bytes(), nil
}

func main() {
	// convertPortrait(buffer Uint8Array) returns a promise that resolves the resulting Uint8Array
	js.Global().Set("convertPortrait", js.FuncOf(func(_ js.Value, args []js.Value) any {
		buf := make([]byte, args[0].Length())
		js.CopyBytesToGo(buf, args[0])
		handler := js.FuncOf(func(_ js.Value, args []js.Value) any {
			resolve := args[0]
			reject := args[1]

			go func() {
				converted, err := convert(buf)
				if err != nil {
					errorObject := js.Global().Get("Error").New(err.Error())
					reject.Invoke(errorObject)
					return
				}
				convertedArray := js.Global().Get("Uint8Array").New(len(converted))
				js.CopyBytesToJS(convertedArray, converted)
				resolve.Invoke(convertedArray)
			}()

			return nil
		})

		return js.Global().Get("Promise").New(handler)
	}))

	// wait forever so wasm function continues to execute
	<-make(chan struct{})
}
