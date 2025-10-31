package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"time"
)

func processImage(imagePath string) (string, error) {
	infile, err := os.Open(imagePath)
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	// Decode the image
	img, _, err := image.Decode(infile)
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	// Iterate through each pixel and convert to grayscale
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			grayColor := color.GrayModel.Convert(originalColor)
			grayImg.Set(x, y, grayColor)
		}
	}

	// --- Ensure temp folder exists ---
	outputDir := "temp"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Println("Failed to create temp folder:", err)
		return "", err
	}

	outputPath, filename := generatePath(outputDir, "window_capture_grayscale_")

	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return "", err
	}
	defer file.Close()

	if err := png.Encode(file, grayImg); err != nil {
		fmt.Println("Failed to encode PNG:", err)
		return "", err
	}

	return filename, nil
}

// func openImage(imagePath string) image.Image {
// 	// Open the input image file

// 	return img
// }

// func createImageFile(imagePath string) (*os.File, error) {

// 	return file, nil
// }

func generatePath(outputDir string, prefix string) (string, string) {

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s%s.png", prefix, timestamp)
	outputPath := filepath.Join(outputDir, filename)

	return outputPath, filename
}
