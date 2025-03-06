package services

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/disintegration/imaging"
	"github.com/starks97/alcohol-tracker-api/utils"
)

func ProcessImage(imgPath io.Reader) ([][]byte, error) {
	//read image and convert to bytes
	imgBytes, err := io.ReadAll(imgPath)
	if err != nil {
		log.Println("Error: Unable to read image")
		return nil, err
	}
	// Load Image
	img, err := imaging.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		log.Println("Error: Unable to decode image")
		return nil, fmt.Errorf("unable to decode image: %v", err)
	}

	// Preprocessing steps
	img = imaging.Grayscale(img)                         // Convert to grayscale
	img = imaging.AdjustContrast(img, 30)                // Increase contrast
	img = imaging.Resize(img, 800, 600, imaging.Lanczos) // Resize image

	denoisedBytes, err := utils.ImageToBytes(img)
	if err != nil {
		log.Println("Error: Unable to convert image to bytes")
		return nil, err
	}

	return [][]byte{denoisedBytes}, nil

}
