package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func analyze(img image.Image) (string, error) {

	// replacer := strings.NewReplacer(
	// 	" ", "-",
	// 	" ", "-",
	// 	"ú", "u",
	// 	":", "",
	// 	"ö", "o",
	// )

	tmpFile, err := ioutil.TempFile("", "ocr-*.png")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Encode image to temp file
	if err := png.Encode(tmpFile, img); err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to encode image: %w", err)
	}
	tmpFile.Close()

	// Call tesseract directly
	cmd := exec.Command("tesseract", tmpPath, "stdout")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("OCR failed:", err)
	}

	text := strings.TrimSpace(string(output))
	fmt.Println("🧠 OCR Result:")
	fmt.Println(text)

	// votingKeywords := []string{"VOTE", "MAP"}
	// for _, keyword := range votingKeywords {
	// 	if strings.Contains(text, keyword) {
	// 		fmt.Printf("Found keyword: %s\n", keyword)
	// 	}
	// }

	return text, nil

	// filePath := "logs/output.txt"
	// timestamp := time.Now().Format("2006-01-02_15-04-05")
	// ocrResult := fmt.Sprintf("🧠 OCR Result:\n%s\n\n%s\n\n", timestamp, text)
	// content := []byte(ocrResult)

	// err = os.WriteFile(filePath, content, 0644) // 0644 sets read/write permissions for owner, read-only for others
	// if err != nil {
	// 	log.Fatalf("Error writing to file: %v", err)
	// }

	// fmt.Println("Content successfully written to output.txt")
}
