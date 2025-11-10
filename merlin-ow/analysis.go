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

	tesseractPath := getTesseractPath()
	// fmt.Println("Using tesseract at:", tesseractPath)
	// fmt.Println("Tesseract exists:", fileExists(tesseractPath))

	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	tessdataPath := filepath.Join(exeDir, "tesseract", "tessdata")
	// log.Println("Executable path:", exePath)
	// log.Println("Executable dir:", exeDir)
	// fmt.Println("Tessdata exists:", fileExists(tessdataPath))

	cmd := exec.Command(tesseractPath, tmpPath, "stdout")
	cmd.Env = append(os.Environ(), "TESSDATA_PREFIX="+tessdataPath)
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("OCR failed:", err)
	}

	text := strings.TrimSpace(string(output))

	return text, nil

}

func getTesseractPath() string {
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, "tesseract", "tesseract.exe")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
