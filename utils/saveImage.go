package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/disintegration/imaging"
)

func SaveImage(imgBytes []byte, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("unable to create file: %v", err)
	}

	defer file.Close()

	img, err := imaging.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return fmt.Errorf("unable to decode image: %v", err)
	}

	// Encode the image to JPEG and save it
	err = imaging.Encode(file, img, imaging.JPEG)
	if err != nil {
		return fmt.Errorf("unable to encode image to JPEG: %v", err)
	}

	log.Printf("File saved successfully as JPEG: %s", fileName)
	return nil

}
