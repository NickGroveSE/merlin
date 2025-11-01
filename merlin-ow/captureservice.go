package main

import (
	"log"
	"os"
	// "time"
)

type CaptureService struct {
	overwatchService *OverwatchService
}

func NewCaptureService(overwatchService *OverwatchService) *CaptureService {
	return &CaptureService{
		overwatchService: overwatchService,
	}
}

func (c *CaptureService) Capture() ([]OWHero, OverwatchFilters, error) {

	captureFilters := OverwatchFilters{Role: "Support", Input: "PC", GameMode: "0", RankTier: "", Map: "new-queen-street", Region: "Americas"}

	// heroes, _ := c.overwatchService.Scrape(captureFilters)

	hero := OWHero{
		Name:     "Ana",
		Color:    "#FF0000",
		Image:    "",
		PickRate: 51.0,
		WinRate:  51.0,
	}

	heroes := []OWHero{hero}

	// for {

	img, err := capture()
	if err != nil {
		log.Fatal("OCR failed:", err)
	}
	processedImg, _ := processImage(img)

	analyze(processedImg)

	// c.removeFile(imagePath)
	// 	time.Sleep(time.Second)
	// }

	return heroes, captureFilters, nil
}

func (c *CaptureService) removeFile(imagePath string) {

	err := os.Remove(imagePath)
	if err != nil {
		log.Fatalf("Error removing file: %v", err)
		return
	}

	log.Printf("File %s removed successfully", imagePath)

}
