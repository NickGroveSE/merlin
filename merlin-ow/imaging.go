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
	"strings"
	"time"

	// contrast "github.com/hawx/img/contrast"
	// sharpen "github.com/hawx/img/sharpen"
	"github.com/andreyvit/locateimage"
)

func imageRecognition(img *image.RGBA) bool {

	selector := RoleSelector{Tank: false, Damage: false, Support: false, Flex: false}

	needleFiles := [8]string{"qp-flex-selected", "qp-tank-selected", "qp-dps-selected", "qp-sup-selected", "comp-flex-selected", "comp-tank-selected", "comp-dps-selected", "comp-sup-selected"}

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
			fmt.Printf("%s Image Not Found", needleFile)
			if strings.Contains(needleFile, "tank") && selector.Tank {
				selector.Tank = false
			} else if strings.Contains(needleFile, "dps") && selector.Damage {
				selector.Damage = false
			} else if strings.Contains(needleFile, "sup") && selector.Support {
				selector.Support = false
			}
		} else if strings.Contains(needleFile, "flex") {
			selector.Tank = false
			selector.Damage = false
			selector.Support = false
			selector.Flex = true
		} else {
			if strings.Contains(needleFile, "tank") && !selector.Tank {
				selector.Tank = true
			} else if strings.Contains(needleFile, "dps") && !selector.Damage {
				selector.Damage = true
			} else if strings.Contains(needleFile, "sup") && !selector.Support {
				selector.Support = true
			}
		}

	}

	selectorReadable := fmt.Sprintf("\nTank Selected: %t\nDamage Selected: %t\nSupport Selected: %t\nFlex Selected: %t", selector.Tank, selector.Damage, selector.Support, selector.Flex)

	fmt.Println(selectorReadable)

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
