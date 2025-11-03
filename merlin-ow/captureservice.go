package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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

func (s Queue) String() string {
	switch s {
	case QP:
		return "Quick Play"
	case Comp:
		return "Competitive"
	case NoSelection:
		return "No Selection"
	default:
		return fmt.Sprintf("Unknown Queue Selection (%d)", s)
	}
}

var qpColor = color.RGBA{R: 27, G: 79, B: 226, A: 255}
var qpColorHover = color.RGBA{R: 28, G: 79, B: 226, A: 255}
var inQueueQPColor = color.RGBA{R: 9, G: 93, B: 222, A: 255}
var compColor = color.RGBA{R: 182, G: 29, B: 71, A: 255}
var compColorHover = color.RGBA{R: 183, G: 30, B: 72, A: 255}
var inQueueCompColor = color.RGBA{R: 208, G: 59, B: 97, A: 255}

// var compColor2 = color.RGBA{R: 161, G: 19, B: 52, A: 255}

func (c *CaptureService) StartMonitoring() ([]OWHero, OverwatchFilters, error) {

	gameState := GameState{GameStatus: StatusIdle, Filters: OverwatchFilters{Input: "PC", Region: "Americas"}, Selector: Selector{}}

	ticker := time.NewTicker(3000 * time.Millisecond)
	defer ticker.Stop()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	counter := 0

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

	c.determineEntryPoint(&gameState)

	for {
		select {
		case <-ticker.C:
			counter++
			go c.evaluate(counter, &gameState)

		case <-sigChan:
			fmt.Println("\nShutting down gracefully...")
			return heroes, captureFilters, nil
		}

	}
	// _ = imageRecognition(img)
	// processedImg, _ := processImage(img)

	// analyze(processedImg)

	// c.removeFile(imagePath)
	// 	time.Sleep(time.Second)
	// }

	return heroes, captureFilters, nil

}

func (c *CaptureService) determineEntryPoint(gameState *GameState) {

	img, err := capture()
	if err != nil {
		log.Fatal("Capture failed:", err)
	}

	queueColorSignifier := img.At(1350, 520)
	inQueueColorSignifier := img.At(1100, 0)
	// fmt.Println(queueColorSignifier)
	// fmt.Println(inQueueColorSignifier)
	switch queueColorSignifier {
	case qpColor, qpColorHover:
		c.updateSelections(img, "qp", QP, gameState)
	case compColor, compColorHover:
		c.updateSelections(img, "comp", Comp, gameState)
	default:
		switch inQueueColorSignifier {
		case inQueueQPColor:
			gameState.GameStatus = StatusInQueue
			gameState.Filters.GameMode = "0"
			fmt.Println(gameState.GameStatus.String())
		case inQueueCompColor:
			gameState.GameStatus = StatusInQueue
			gameState.Filters.GameMode = "1"
			fmt.Println(gameState.GameStatus.String())
		default:
			processedImg, _ := processImage(img)
			text, err := analyze(processedImg)
			if err != nil {
				fmt.Printf("Error in Text Analysis: %e", err)
			}

			if strings.Contains(text, "VOTE") && strings.Contains(text, "MAP") {
				gameState.GameStatus = StatusMapVotingPhase
				fmt.Println(gameState.GameStatus.String())
			} else if strings.Contains(text, "BANNING") {
				gameState.GameStatus = StatusBanningPhase
				fmt.Println(gameState.GameStatus.String())
				// Transmit to User That Starting Monitoring in Banning Phase is Not Recommended
			} else {
				fmt.Println(gameState.GameStatus.String())
			}
		}

	}
}

func (c *CaptureService) evaluate(counter int, gameState *GameState) {
	img, err := capture()
	if err != nil {
		log.Fatal("OCR failed:", err)
	}
	fmt.Printf("Capture %d\n", counter)
	queueColorSignifier := img.At(1350, 520)
	inQueueColorSignifier := img.At(1100, 0)
	// fmt.Println(queueColorSignifier)
	// fmt.Println(inQueueColorSignifier)
	switch queueColorSignifier {
	case qpColor, qpColorHover:
		c.updateSelections(img, "qp", QP, gameState)
	case compColor, compColorHover:
		c.updateSelections(img, "comp", Comp, gameState)
	default:
		switch inQueueColorSignifier {
		case inQueueQPColor:
			c.confirmSelections(img, "0", gameState)
		case inQueueCompColor:
			c.confirmSelections(img, "1", gameState)
		default:
			processedImg, _ := processImage(img)
			text, err := analyze(processedImg)
			if err != nil {
				fmt.Printf("Error in Text Analysis: %e", err)
			}

			if strings.Contains(text, "VOTE") && strings.Contains(text, "MAP") {
				gameState.GameStatus = StatusMapVotingPhase
				fmt.Println(gameState.GameStatus.String())
			} else if strings.Contains(text, "BANNING") {
				gameState.GameStatus = StatusBanningPhase
				fmt.Println(gameState.GameStatus.String())
				// Transmit to User That Starting Monitoring in Banning Phase is Not Recommended
			} else {
				fmt.Println(gameState.GameStatus.String())
			}
		}
	}
}

func (c *CaptureService) updateSelections(img *image.RGBA, queue string, queueEnum Queue, gameState *GameState) {
	gameState.Selector.Queue = queueEnum
	gameState.GameStatus = StatusSelection
	fmt.Println(gameState.GameStatus.String())
	fmt.Println(gameState.Selector.Queue.String())
	_ = roleImageRecognition(img, queue, gameState)
}

func (c *CaptureService) confirmSelections(img *image.RGBA, queueQueryParam string, gameState *GameState) {
	gameState.GameStatus = StatusInQueue
	gameState.Filters.GameMode = queueQueryParam
	if gameState.Selector.Flex {
		fmt.Println("Flex Selection")
		// Transmit Message About Flex
	} else if gameState.Selector.Tank {
		if gameState.Selector.Damage {
			fmt.Println("Multiple Role Selection")
			gameState.Filters.Role = "Tank"
			// Transmit Message About Multiple Role Selection
		} else if gameState.Selector.Support {
			fmt.Println("Multiple Role Selection")
			gameState.Filters.Role = "Tank"
			// Transmit Message About Multiple Role Selection
		} else {
			fmt.Println("Tank")
			gameState.Filters.Role = "Tank"
		}
	} else if gameState.Selector.Damage {
		if gameState.Selector.Support {
			fmt.Println("Multiple Role Selection")
			gameState.Filters.Role = "Damage"
			// Transmit Message About Multiple Role Selection
		} else {
			fmt.Println("Damage")
			gameState.Filters.Role = "Damage"
		}
	} else if gameState.Selector.Support {
		fmt.Println("Support")
		gameState.Filters.Role = "Support"
	} else {
		fmt.Println("Role Selection Couldn't Be Detected")
		gameState.Filters.Role = "Tank"
		// Transmit Message About Role Selection Couldn't Be Detected
	}

	gameState.Selector = Selector{}

	fmt.Println(gameState.Filters)
}

func (c *CaptureService) removeFile(imagePath string) {

	err := os.Remove(imagePath)
	if err != nil {
		log.Fatalf("Error removing file: %v", err)
		return
	}

	log.Printf("File %s removed successfully", imagePath)

}
