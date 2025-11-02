package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"time"

	// contrast "github.com/hawx/img/contrast"
	// sharpen "github.com/hawx/img/sharpen"
	"github.com/andreyvit/locateimage"
)

type Selection struct {
	Queue   Queue
	Tank    bool
	Damage  bool
	Support bool
	Flex    bool
}

type Queue int

const (
	QP Queue = iota
	Comp
	NoSelection
)

func imageRecognition(img *image.RGBA) bool {

	// selection := Selection{Queue: QP, Tank: false, Damage: false, Support: false, Flex: false}

	needleFiles := [7]string{"qp-flex-selected", "qp-tank-selected", "qp-tank-unselected", "qp-dps-selected", "qp-dps-unselected", "qp-sup-selected", "qp-sup-unselected"}

	for _, needleFile := range needleFiles {

		file, err := os.Open("recognition_assets/" + needleFile + ".png")
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return false
		}
		defer file.Close() // Ensure the file is closed when done

		// Decode the PNG image
		needle, err := png.Decode(file)
		if err != nil {
			fmt.Printf("Error decoding PNG: %v\n", err)
			return false
		}

		needleRGBA := imageToRGBA(needle)

		_, err = locateimage.Find(context.Background(), img, needleRGBA, 0, locateimage.Fastest)
		if err != nil {
			fmt.Printf("Error Finding Image: %v\n", err)
		} else {
			fmt.Println("Holy Shit It Worked")
		}

	}

	return true
}

func processImage(img *image.RGBA) (image.Image, error) {

	// processedImg := sharpen.UnsharpMask(contrast.Adjust(grayscale(img), 20), 3, 2.5, 2.5, 0.0)
	// processedImg := contrast.Adjust(grayscale(img), 20)
	processedImg := grayscale(img)

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

func imageToRGBA(src image.Image) *image.RGBA {
	// If the source is already an *image.RGBA, no conversion is needed.
	if dst, ok := src.(*image.RGBA); ok {
		return dst
	}

	// Create a new RGBA image with the same bounds as the source image.
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))

	// Draw the source image onto the new RGBA image.
	// This performs the necessary color model conversion.
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}
