package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

	// Then replace your existing code with:
	tesseractPath := getTesseractPath()

	// Set environment so tesseract finds its data files
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	os.Setenv("TESSDATA_PREFIX", filepath.Join(exeDir, "tesseract"))

	// Call tesseract directly
	cmd := exec.Command(tesseractPath, tmpPath, "stdout")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("OCR failed:", err)
	}

	text := strings.TrimSpace(string(output))
	// fmt.Println("🧠 OCR Result:")
	// fmt.Println(text)

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

func getTesseractPath() string {
	exePath, err := os.Executable()
	if err != nil {
		return "tesseract" // fallback for dev mode
	}

	exeDir := filepath.Dir(exePath)
	tesseractPath := filepath.Join(exeDir, "tesseract", "tesseract.exe")

	// Check if bundled version exists
	if _, err := os.Stat(tesseractPath); err == nil {
		return tesseractPath // use bundled version
	}

	return "tesseract" // fallback to system version
}
