package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	// "time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type CaptureService struct {
	overwatchService *OverwatchService
	app              *application.App
	stopChan         chan struct{}
	isRunning        bool
	mu               sync.Mutex
}

func NewCaptureService(overwatchService *OverwatchService) *CaptureService {
	return &CaptureService{
		overwatchService: overwatchService,
		app:              nil,
		stopChan:         nil,
		isRunning:        false,
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

// type StatusUpdate struct {
// 	StatusIcon string `json:"statusIcon"`
// 	StatusText string `json:"statusText"`
// 	Message    string `json:"message"`
// }

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

var mapScan = [31]string{
	"HANOAKA",
	"THRONE OF ANUBIS",
	"ANTARCTIC PENINSULA",
	"BUSAN",
	"ILIOS",
	"LIGJANG TOWER",
	"NEPAL",
	"OASIS",
	"SAMOA",
	"CIRCUIT ROYAL",
	"DORADO",
	"HAVANA",
	"JUNKERTOWN",
	"RIALTO",
	"ROUTE 66",
	"SHAMBALI MONASTERY",
	"WATCHPOINT: GIBRALTAR",
	"AATLIS",
	"NEW JUNK CITY",
	"SURVASA",
	"BLIZZARD WORLD",
	"EICHENWALDE",
	"HOLLYWOOD",
	"KING'S ROW",
	"MIDTOWN",
	"NUMBANI",
	"PARAÍSO",
	"COLOSSEO",
	"ESPERANÇA",
	"NEW QUEEN STREET",
	"RUNASAPI",
}

var mapFormat = map[string]string{
	"HANAOKA":               "hanaoka",
	"THRONE OF ANUBIS":      "throne-of-anubis",
	"ANTARCTIC PENINSULA":   "antarctic-peninsula",
	"BUSAN":                 "busan",
	"ILIOS":                 "ilios",
	"LIJANG TOWER":          "lijang-tower",
	"NEPAL":                 "nepal",
	"OASIS":                 "oasis",
	"SAMOA":                 "samoa",
	"CIRCUIT ROYAL":         "circuit-royal",
	"DORADO":                "dorado",
	"HAVANA":                "havana",
	"JUNKERTOWN":            "junkertown",
	"RIALTO":                "rialto",
	"ROUTE 66":              "route-66",
	"SHAMBALI MONASTERY":    "shambali-monastery",
	"WATCHPOINT: GIBRALTAR": "watchpoint-gibraltar",
	"AATLIS":                "aatlis",
	"NEW JUNK CITY":         "new-junk-city",
	"SURVASA":               "survasa",
	"BLIZZARD WORLD":        "blizzard-world",
	"EICHENWALDE":           "eichenwalde",
	"HOLLYWOOD":             "hollywood",
	"KING'S ROW":            "kings-row",
	"MIDTOWN":               "midtown",
	"NUMBANI":               "numbani",
	"PARAÍSO":               "paraiso",
	"COLOSSEO":              "colosseo",
	"ESPERANÇA":             "esperanca",
	"NEW QUEEN STREET":      "new-queen-street",
	"RUNASAPI":              "runasapi",
}

var defaultStatusMessage = "Messages with more in-depth status updates..."

func (c *CaptureService) StartMonitoring() ([]OWHero, OverwatchFilters, error) {

	c.mu.Lock()
	if c.isRunning {
		c.mu.Unlock()
		return []OWHero{}, OverwatchFilters{}, errors.New("monitoring already running")
	}
	c.isRunning = true
	c.stopChan = make(chan struct{})
	c.mu.Unlock()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	doneChan := make(chan struct{})

	gameState := GameState{GameStatus: StatusIdle, Filters: OverwatchFilters{Role: "Support", Input: "PC", GameMode: "2", RankTier: "All", Map: "all-maps", Region: "Americas"}, Selector: Selector{}}

	counter := 0

	// captureFilters := OverwatchFilters{Role: "Support", Input: "PC", GameMode: "0", RankTier: "", Map: "new-queen-street", Region: "Americas"}

	// heroes, _ := c.overwatchService.Scrape(captureFilters)

	// hero := OWHero{
	// 	Name:     "Ana",
	// 	Color:    "#FF0000",
	// 	Image:    "",
	// 	PickRate: 51.0,
	// 	WinRate:  51.0,
	// }

	// heroes := []OWHero{hero}

	c.determineEntryPoint(&gameState)

	for {
		select {
		case <-ticker.C:
			counter++
			go c.evaluate(counter, &gameState, doneChan)

		case <-doneChan:
			fmt.Println("\nFinished retrieving filters, Scraping for data...")
			heroes, err := c.overwatchService.Scrape(gameState.Filters)
			if err != nil {
				fmt.Printf("Error Scraping: %e", err)
			}
			return heroes, gameState.Filters, nil

		case <-c.stopChan:
			fmt.Println("\nMonitoring cut early, Scraping for data...")
			heroes, err := c.overwatchService.Scrape(gameState.Filters)
			if err != nil {
				fmt.Printf("Error Scraping: %e", err)
			}
			return heroes, gameState.Filters, nil

		case <-sigChan:
			fmt.Println("\nShutting down gracefully...")
			return []OWHero{}, OverwatchFilters{}, nil
		}

	}
	// _ = imageRecognition(img)
	// processedImg, _ := processImage(img)

	// analyze(processedImg)

	// c.removeFile(imagePath)
	// 	time.Sleep(time.Second)
	// }

	return []OWHero{}, OverwatchFilters{}, nil

}

func (c *CaptureService) StopMonitoring() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isRunning && c.stopChan != nil {
		c.stopChan <- struct{}{}
		c.isRunning = false
	}
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
			gameState.Filters.GameMode = "2"
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
			} else if strings.Contains(text, "SELECT") {
				gameState.GameStatus = StatusBanningPhase
				fmt.Println(gameState.GameStatus.String())
				// Transmit to User That Starting Monitoring in Banning Phase is Not Recommended
			} else {
				fmt.Println(gameState.GameStatus.String())
			}
		}

	}
}

func (c *CaptureService) evaluate(counter int, gameState *GameState, done chan struct{}) {

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
		if gameState.GameStatus != StatusSelection {
			c.app.Event.Emit("status-update", map[string]string{
				"statusIcon": "../public/assets/selection.svg",
				"statusText": "Selecting Role",
				"message":    "Select your role(s) and game mode and I'll lock them in once you queue.",
			})
			c.app.Event.Emit("queue-update", "Waiting...")
			c.app.Event.Emit("role-update", "Waiting...")
			c.updateSelections(img, "qp", QP, gameState)
			// c.app.Event.Emit("test-emit")
		}
	case compColor, compColorHover:
		if gameState.GameStatus != StatusSelection {
			c.app.Event.Emit("status-update", map[string]string{
				"statusIcon": "../public/assets/selection.svg",
				"statusText": "Selecting Role",
				"message":    "Select your role(s) and game mode and I'll will lock them in once you queue.",
			})
		}
		c.app.Event.Emit("queue-update", "Waiting...")
		c.app.Event.Emit("role-update", "Waiting...")
		c.updateSelections(img, "comp", Comp, gameState)
	default:
		if gameState.GameStatus == StatusSelection || gameState.GameStatus == StatusIdle {
			switch inQueueColorSignifier {
			case inQueueQPColor:
				if gameState.GameStatus != StatusInQueue {
					c.app.Event.Emit("status-update", map[string]string{
						"statusIcon": "../public/assets/in-queue.svg",
						"statusText": "In Queue",
						"message":    "",
					})
					c.app.Event.Emit("queue-update", "Quick Play")
					c.app.Event.Emit("map-update", "Waiting for Match...")
					c.confirmSelections("0", gameState)
				}
			case inQueueCompColor:
				if gameState.GameStatus != StatusInQueue {
					c.app.Event.Emit("status-update", map[string]string{
						"statusIcon": "../public/assets/in-queue.svg",
						"statusText": "In Queue",
						"message":    "",
					})
					c.app.Event.Emit("queue-update", "Competitive")
					c.app.Event.Emit("map-update", "Waiting for Match...")
					c.confirmSelections("2", gameState)
				}
			default:
				fmt.Println(gameState.GameStatus.String())
			}
		} else if gameState.GameStatus == StatusInQueue {
			processedImg, _ := processImage(img)
			text, err := analyze(processedImg)
			if err != nil {
				fmt.Printf("Error in Text Analysis: %e", err)
			}

			if strings.Contains(text, "VOTE") && strings.Contains(text, "MAP") {
				if gameState.GameStatus != StatusMapVotingPhase {
					c.app.Event.Emit("status-update", map[string]string{
						"statusIcon": "../public/assets/map-voting.svg",
						"statusText": "Map Voting",
						"message":    "I see you've entered a match! I'll collect our final data needed after the map vote",
					})
				}
				gameState.GameStatus = StatusMapVotingPhase
				fmt.Println(gameState.GameStatus.String())
			} else {
				fmt.Println(gameState.GameStatus.String())
			}
		} else if gameState.GameStatus == StatusMapVotingPhase {
			processedImg, _ := processImage(img)
			text, err := analyze(processedImg)
			if err != nil {
				fmt.Printf("Error in Text Analysis: %e", err)
			}

			fmt.Println(text)

			if strings.Contains(text, "RESULT") {
				if gameState.GameStatus != StatusBanningPhase {
					c.app.Event.Emit("status-update", map[string]string{
						"statusIcon": "../public/assets/banning-phase.svg",
						"statusText": "Bans Phase",
						"message":    "Detecting map vote result...",
					})
				}

				gameState.GameStatus = StatusBanningPhase
				postVoteText, err := analyze(processedImg)
				if err != nil {
					fmt.Printf("Error in Text Analysis: %e", err)
				}

				for i := range len(mapScan) {
					if strings.Contains(postVoteText, mapScan[i]) {
						gameState.Filters.Map = mapFormat[mapScan[i]]
						done <- struct{}{}
					}
				}
				fmt.Println(gameState.GameStatus.String())
			} else {
				fmt.Println(gameState.GameStatus.String())
			}
		} else if gameState.GameStatus == StatusBanningPhase {
			mapCap, err := captureMap()
			if err != nil {
				log.Fatal("Capture failed:", err)
			}

			processedMapCap, _ := processImage(mapCap)
			mapText, err := analyze(processedMapCap)
			if err != nil {
				fmt.Printf("Error in Text Analysis: %e", err)
			}

			for i := range len(mapScan) {
				if strings.Contains(mapText, mapScan[i]) {
					gameState.Filters.Map = mapFormat[mapScan[i]]
					done <- struct{}{}
				}
			}
		} else {
			fmt.Println(gameState.GameStatus.String())
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

func (c *CaptureService) confirmSelections(queueQueryParam string, gameState *GameState) {
	switch gameState.GameStatus {
	case StatusSelection:
		gameState.GameStatus = StatusInQueue
		gameState.Filters.GameMode = queueQueryParam
		if gameState.Selector.Flex {
			fmt.Println("Flex Selection")
			gameState.Filters.Role = "Tank"
			c.app.Event.Emit("message", "Flex Queue! A jack of all trades I see. I'll put you as Tank for now but once we have collected the match data for you, you'll be able to pick what role you end up with in the filters")
			c.app.Event.Emit("role-update", "Tank")
		} else if gameState.Selector.Tank {
			if gameState.Selector.Damage {
				fmt.Println("Multiple Role Selection")
				gameState.Filters.Role = "Tank"
				c.app.Event.Emit("message", "I see you have chosen multiple roles, I'll put you as Tank for now but once we have collected the match data for you, you'll be able to pick what role you end up with in the filters")
				c.app.Event.Emit("role-update", "Tank")
			} else if gameState.Selector.Support {
				fmt.Println("Multiple Role Selection")
				gameState.Filters.Role = "Tank"
				c.app.Event.Emit("message", "I see you have chosen multiple roles, I'll put you as Tank for now but once we have collected the match data for you, you'll be able to pick what role you end up with in the filters")
				c.app.Event.Emit("role-update", "Tank")
			} else {
				fmt.Println("Tank")
				gameState.Filters.Role = "Tank"
				c.app.Event.Emit("message", "Locked in your role as Tank")
				c.app.Event.Emit("role-update", "Tank")
			}
		} else if gameState.Selector.Damage {
			if gameState.Selector.Support {
				fmt.Println("Multiple Role Selection")
				gameState.Filters.Role = "Damage"
				c.app.Event.Emit("message", "I see you have chosen multiple roles, I'll put you as Damage for now but once we have collected the match data for you, you'll be able to pick what role you end up with in the filters")
				c.app.Event.Emit("role-update", "Damage")
			} else {
				fmt.Println("Damage")
				gameState.Filters.Role = "Damage"
				c.app.Event.Emit("message", "Locked in your role as Damage")
				c.app.Event.Emit("role-update", "Damage")
			}
		} else if gameState.Selector.Support {
			fmt.Println("Support")
			gameState.Filters.Role = "Support"
			c.app.Event.Emit("message", "Locked in your role as Support")
			c.app.Event.Emit("role-update", "Support")
		} else {
			fmt.Println("Role Selection Couldn't Be Detected")
			gameState.Filters.Role = "Tank"
			c.app.Event.Emit("message", "Oops! Looks like I couldn't detect your role. Its okay though once we have collected the other match data for you, you'll be able to pick what role you end up with in the filters")
			c.app.Event.Emit("role-update", "Support")
		}
	case StatusIdle:
		gameState.GameStatus = StatusInQueue
		gameState.Filters.GameMode = queueQueryParam
	}

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
