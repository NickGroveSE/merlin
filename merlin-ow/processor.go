package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	// "image/png"
	// "os"
	"path/filepath"
	"time"
	// contrast "github.com/hawx/img/contrast"
	// sharpen "github.com/hawx/img/sharpen"
	// "github.com/andreyvit/locateimage"
)

func roleRecognition(cap Capture, gameState *GameState) error {

	xBoundCalcs := [4]int{
		int(cap.xBoundForCalc * 0.7197916666666667),
		int(cap.xBoundForCalc * 0.559375),
		int(cap.xBoundForCalc * 0.4109375),
		int(cap.xBoundForCalc * 0.2625),
	}

	checkColor := color.RGBA{R: 0, G: 186, B: 0, A: 255}
	consolidatedYBoundCalc := int(cap.yBoundForCalc * 0.5185185185185185)
	// uncheckColor := color.RGBA{R: 29, G: 37, B: 58, A: 255}

	for i, xBoundCalc := range xBoundCalcs {

		// pixel := img.At(coord, 587)

		// if ColorMatch(uncheckColor, pixel, 2000) {
		// 	if i == 0 && gameState.Selector.Flex {
		// 		fmt.Println("Flex Unselected")
		// 		gameState.Selector.Flex = false
		// 	} else if i == 1 && gameState.Selector.Support {
		// 		fmt.Println("Support Unselected")
		// 		gameState.Selector.Support = false
		// 	} else if i == 2 && gameState.Selector.Damage {
		// 		fmt.Println("Damage Unselected")
		// 		gameState.Selector.Damage = false
		// 	} else if i == 3 && gameState.Selector.Tank {
		// 		fmt.Println("Tank Unselected")
		// 		gameState.Selector.Tank = false
		// 	}
		// } else
		if ColorExistsInRegion(cap.img, xBoundCalc, consolidatedYBoundCalc, 44, 44, checkColor) {
			if i == 0 && !gameState.Selector.Flex {
				fmt.Println("Flex Selected")
				gameState.Selector.Flex = true
				gameState.Selector.Support = false
				gameState.Selector.Damage = false
				gameState.Selector.Tank = false
			} else if i == 1 && !gameState.Selector.Support {
				fmt.Println("Support Selected")
				gameState.Selector.Support = true
			} else if i == 2 && !gameState.Selector.Damage {
				fmt.Println("Damage Selected")
				gameState.Selector.Damage = true
			} else if i == 3 && !gameState.Selector.Tank {
				fmt.Println("Tank Selected")
				gameState.Selector.Tank = true
			}
		} else {
			if i == 0 && gameState.Selector.Flex {
				fmt.Println("Flex Unselected")
				gameState.Selector.Flex = false
			} else if i == 1 && gameState.Selector.Support {
				fmt.Println("Support Unselected")
				gameState.Selector.Support = false
			} else if i == 2 && gameState.Selector.Damage {
				fmt.Println("Damage Unselected")
				gameState.Selector.Damage = false
			} else if i == 3 && gameState.Selector.Tank {
				fmt.Println("Tank Unselected")
				gameState.Selector.Tank = false
			}
		}

	}

	return nil
}

func ColorExistsInRegion(img *image.RGBA, x, y, width, height int, targetColor color.Color) bool {
	tr, tg, tb, ta := targetColor.RGBA()

	for py := y; py < y+height; py++ {
		for px := x; px < x+width; px++ {
			r, g, b, a := img.At(px, py).RGBA()
			if r == tr && g == tg && b == tb && a == ta {
				return true
			}
		}
	}
	return false
}

func ColorMatch(c1 color.Color, c2 color.Color, threshold uint32) bool {

	// fmt.Println(c1)
	// fmt.Println(c2)

	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()

	totalDiff := Abs(int32(r1-r2)) + Abs(int32(g1-g2)) + Abs(int32(b1-b2))

	// fmt.Println(totalDiff)

	return totalDiff < int32(threshold)
}

func Abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
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
