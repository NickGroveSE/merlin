package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	imagePath := "temp/" + captureHandler()

	// Call tesseract directly
	cmd := exec.Command("tesseract", imagePath, "stdout")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("OCR failed:", err)
	}

	text := strings.TrimSpace(string(output))
	fmt.Println("🧠 OCR Result:")
	fmt.Println(text)
}
