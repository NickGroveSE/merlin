package main

import (
	"fmt"
	"image/color"
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
	StatusSelection
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
	case StatusSelection:
		return "Selecting Queue & Role"
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

var qpColor = color.RGBA{R: 13, G: 56, B: 217, A: 255}
var compColor = color.RGBA{R: 161, G: 20, B: 51, A: 255}

func (c *CaptureService) StartMonitoring() ([]OWHero, OverwatchFilters, error) {

	gameState := GameState{GameStatus: StatusIdle, Filters: OverwatchFilters{}, Selector: Selector{}}

	// for {

	img, err := capture()
	if err != nil {
		log.Fatal("OCR failed:", err)
	}
	queueColorSignifier := img.At(500, 600)
	switch queueColorSignifier {
	case qpColor:
		gameState.Selector.Queue = QP
		gameState.GameStatus = StatusSelection
		_ = imageRecognition(img, "qp", &gameState)
		fmt.Println("We Are Selecting for QP")
	case compColor:
		gameState.Selector.Queue = Comp
		gameState.GameStatus = StatusSelection
		_ = imageRecognition(img, "qp", &gameState)
		fmt.Println("We Are Selecting for Comp")
	default:
		gameState.Selector.Queue = Comp
		gameState.GameStatus = StatusSelection
		_ = imageRecognition(img, "qp", &gameState)
		fmt.Println("Not Selecting")
	}
	// _ = imageRecognition(img)
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
