package main

import (
	"fmt"
	"log"
	"os"

	// "time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type CaptureService struct {
	overwatchService *OverwatchService
	app              *application.App
}

func NewCaptureService(overwatchService *OverwatchService) *CaptureService {
	return &CaptureService{
		overwatchService: overwatchService,
		app:              nil,
	}
}

func (c *CaptureService) SetApp(app *application.App) {
	c.app = app
}

type GameState struct {
	GameStatus Status
	Filters    OverwatchFilters
	Selector   Selector
}

type Status int

const (
	StatusIdle Status = iota
	StatusRoleSelection
	StatusInQueue
	StatusMapVotingPhase
	StatusBanningPhase
)

type Selector struct {
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

func (s Status) String() string {
	switch s {
	case StatusIdle:
		return "Idle"
	case StatusRoleSelection:
		return "Selecting Role"
	case StatusInQueue:
		return "In Queue"
	case StatusMapVotingPhase:
		return "Voting for a Map"
	case StatusBanningPhase:
		return "Banning Phase"
	default:
		return fmt.Sprintf("Unknown Status (%d)", s)
	}
}

func (c *CaptureService) StartMonitoring() ([]OWHero, OverwatchFilters, error) {

	// gameState := GameState{GameStatus: StatusIdle, Filters: OverwatchFilters{}}

	// for {

	img, err := capture()
	if err != nil {
		log.Fatal("OCR failed:", err)
	}
	_ = imageRecognition(img)
	// processedImg, _ := processImage(img)

	// analyze(processedImg)

	// c.removeFile(imagePath)
	// 	time.Sleep(time.Second)
	// }

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
