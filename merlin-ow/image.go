package main

import (
	"fmt"
	"image"
	"image/color"
	"path/filepath"
	"time"

	contrast "github.com/hawx/img/contrast"
	// sharpen "github.com/hawx/img/sharpen"
)

func processImage(img *image.RGBA) (image.Image, error) {

	// processedImg := sharpen.UnsharpMask(contrast.Adjust(grayscale(img), 20), 3, 2.5, 2.5, 0.0)
	processedImg := contrast.Adjust(grayscale(img), 20)

	return processedImg, nil
}

func grayscale(img image.Image) *image.Gray {

	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	if rgba, ok := img.(*image.RGBA); ok {
		// Fast memory copy for RGBA
		idx := 0
		for i := 0; i < len(rgba.Pix); i += 4 {
			r, g, b := rgba.Pix[i], rgba.Pix[i+1], rgba.Pix[i+2]
			grayImg.Pix[idx] = uint8((19595*uint32(r) + 38470*uint32(g) + 7471*uint32(b) + 1<<15) >> 16)
			idx++
		}
		return grayImg
	}

	// Fallback
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayImg.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
		}
	}
	return grayImg

}

func generatePath(outputDir string, prefix string) (string, string) {

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s%s.png", prefix, timestamp)
	outputPath := filepath.Join(outputDir, filename)

	return outputPath, filename
}
