package main

import (
	"fmt"
	"log"

	// "os"
	"os/exec"
	"strings"
	// "time"
)

func analyze(imagePath string) {

	// replacer := strings.NewReplacer(
	// 	" ", "-",
	// 	" ", "-",
	// 	"ú", "u",
	// 	":", "",
	// 	"ö", "o",
	// )

	// Call tesseract directly
	cmd := exec.Command("tesseract", imagePath, "stdout")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("OCR failed:", err)
	}

	text := strings.TrimSpace(string(output))
	fmt.Println("🧠 OCR Result:")
	fmt.Println(text)

	votingKeywords := []string{"VOTE", "MAP"}

	for _, keyword := range votingKeywords {
		if strings.Contains(text, keyword) {
			fmt.Printf("%s\n", keyword)
		}
	}

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
