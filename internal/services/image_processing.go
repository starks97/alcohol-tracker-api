package services

import (
	"fmt"
	"image"
	"io"
	"log"

	"gocv.io/x/gocv"
)

func ProcessImage(imgPath io.Reader) ([]byte, error) {
	//read image and convert to bytes
	imgBytes, err := io.ReadAll(imgPath)
	if err != nil {
		log.Println("Error: Unable to read image")
		return nil, err
	}
	// Load Image
	img, err := gocv.IMDecode(imgBytes, gocv.IMReadColor)
	if err != nil || img.Empty() {
		log.Println("Error: Unable to decode image")
		return nil, fmt.Errorf("unable to decode image")
	}
	defer img.Close()

	// Convert to Grayscale
	gray := gocv.NewMat()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)
	defer gray.Close()

	// Apply Bilateral Filtering (denoises but keeps edges)
	denoised := gocv.NewMat()
	gocv.BilateralFilter(gray, &denoised, 9, 75, 75)
	defer denoised.Close()

	// Apply CLAHE (Contrast Limited Adaptive Histogram Equalization)
	clahe := gocv.NewCLAHEWithParams(2.0, image.Pt(8, 8)) // Fixed
	defer clahe.Close()

	enhanced := gocv.NewMat()
	clahe.Apply(denoised, &enhanced)
	defer enhanced.Close()

	// Apply Thresholding (removes background)
	thresh := gocv.NewMat()
	gocv.Threshold(enhanced, &thresh, 100, 255, gocv.ThresholdBinary+gocv.ThresholdOtsu)
	defer thresh.Close()

	// Apply Morphological Operations (removes small noise)
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
	defer kernel.Close()

	morph := gocv.NewMat()
	gocv.MorphologyEx(thresh, &morph, gocv.MorphClose, kernel)
	defer morph.Close()

	// Sharpen Image using Laplacian filter
	laplacian := gocv.NewMat()
	gocv.Laplacian(morph, &laplacian, gocv.MatTypeCV8U, 3, 1, 0, gocv.BorderDefault)
	defer laplacian.Close()

	// Normalize Image (scaling between 0-255)
	normalized := gocv.NewMat()
	gocv.Normalize(laplacian, &normalized, 0, 255, gocv.NormMinMax)
	defer normalized.Close()

	// Resize to match NN input size (Example: 224x224)
	resized := gocv.NewMat()
	gocv.Resize(normalized, &resized, image.Pt(224, 224), 0, 0, gocv.InterpolationLinear)

	buf, err := gocv.IMEncode(gocv.JPEGFileExt, resized)
	if err != nil {
		log.Printf("Error encoding the image: %v", err)
		return nil, err
	}
	defer buf.Close()

	return buf.GetBytes(), nil // Clone ensures returned Mat is not deallocated
}
