package utils

import (
	"bytes"
	"fmt"
	"image"

	"github.com/disintegration/imaging"
)

func ImageToBytes(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := imaging.Encode(&buf, img, imaging.JPEG)
	if err != nil {
		return nil, fmt.Errorf("unable to encode image to bytes: %v", err)
	}
	return buf.Bytes(), nil
}
