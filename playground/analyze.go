package main

import (
	"fmt"
	"image/color"
	"time"
)

func main() {
	// imagePath := "temp/" + captureHandler()

	// Call tesseract directly
	// cmd := exec.Command("tesseract", imagePath, "stdout")
	// output, err := cmd.Output()
	// if err != nil {
	// 	log.Fatal("OCR failed:", err)
	// }

	// text := strings.TrimSpace(string(output))
	// fmt.Println("🧠 OCR Result:")
	// fmt.Println(text)

	img := captureHandler()

	// cip := img.At(1133, 535)
	// fmt.Println(cip) // (182, 30, 72)
	// qip := img.At(1133, 535)
	// fmt.Println(qip) // (25 79 227)
	// fcp := img.At(1394, 587)
	// fmt.Println(fcp) // (0 186 0)
	// scp := img.At(1086, 587)
	// fmt.Println(scp) // (0 186 0)
	// dcp := img.At(801, 587)
	// fmt.Println(dcp) // (0 186 0)
	tcp := img.At(516, 587)
	fmt.Println(tcp) // (0 186 0)

	// Not Selected RGB (29 37 58)

	time.Sleep(3 * time.Second)

	img = captureHandler()

	// cihp := img.At(1133, 535)
	// fmt.Println(cihp) // (182, 28, 73)
	// qihp := img.At(1133, 535)
	// fmt.Println(qihp) // (14 56 216)
	// fchp := img.At(1394, 587)
	// fmt.Println(fchp) // (9 167 24)
	// schp := img.At(1086, 587)
	// fmt.Println(schp) // (9 167 24)
	// dchp := img.At(801, 587)
	// fmt.Println(dchp) // (9 167 24)
	tchp := img.At(516, 587)
	fmt.Println(tchp) // (9 167 24)

	// fmt.Println(colorMatch(cip, cihp, 15000))
	// fmt.Println(colorMatch(qip, qihp, 15000))
	// fmt.Println(colorMatch(fcp, fchp, 15000))
	// fmt.Println(colorMatch(scp, schp, 15000))
	// fmt.Println(colorMatch(dcp, dchp, 15000))
	fmt.Println(colorMatch(tcp, tchp, 15000))

}

func colorMatch(c1 color.Color, c2 color.Color, threshold uint32) bool {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()

	totalDiff := abs(int32(r1-r2)) + abs(int32(g1-g2)) + abs(int32(b1-b2))

	fmt.Println(totalDiff)

	return totalDiff < int32(threshold)
}

func abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
