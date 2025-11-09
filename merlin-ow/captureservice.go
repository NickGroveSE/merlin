package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"

	// "log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

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

type Capture struct {
	img           *image.RGBA
	xBound        int32
	yBound        int32
	xBoundForCalc float64
	yBoundForCalc float64
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

var qpColor = color.RGBA{R: 25, G: 79, B: 227, A: 255}

// var qpColorHover = color.RGBA{R: 28, G: 79, B: 226, A: 255}
var inQueueQPColor = color.RGBA{R: 9, G: 93, B: 222, A: 255}
var compColor = color.RGBA{R: 182, G: 30, B: 72, A: 255}

// var compColorHover = color.RGBA{R: 183, G: 30, B: 72, A: 255}
var inQueueCompColor = color.RGBA{R: 208, G: 59, B: 97, A: 255}

var mapScan = [31]string{
	"HANOAKA",
	"THRONE OF ANUBIS",
	"ANTARCTIC PENINSULA",
	"BUSAN",
	"ILIOS",
	"LIJIANG TOWER",
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
	"SURAVASA",
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
	"SURAVASA":              "suravasa",
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

	// Ensure cleanup happens
	defer func() {
		c.mu.Lock()
		c.isRunning = false
		c.stopChan = nil
		c.mu.Unlock()
	}()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	doneChan := make(chan struct{}, 1)

	gameState := GameState{
		GameStatus: StatusIdle,
		Filters: OverwatchFilters{
			Role:     "Support",
			Input:    "PC",
			GameMode: "2",
			RankTier: "All",
			Map:      "all-maps",
			Region:   "Americas"},
		Selector: Selector{},
	}

	counter := 0

	if !c.determineEntryPoint(&gameState) {
		return []OWHero{}, OverwatchFilters{}, nil
	}

	for {
		select {
		case <-ticker.C:
			counter++
			c.evaluate(counter, &gameState, doneChan)

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

	return []OWHero{}, OverwatchFilters{}, nil

}

func (c *CaptureService) StopMonitoring() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isRunning && c.stopChan != nil {
		c.app.Event.Emit("message", "Shutting down monitoring software, this will take a moment...")
		close(c.stopChan)
	}
}
func (c *CaptureService) determineEntryPoint(gameState *GameState) bool {

	cap, err := capture()
	if err != nil {
		c.app.Event.Emit("error", err.Error()+"\n\nAn error occured while monitoring your screen, make sure Overwatch is running and set to 'Borderless Windowed' in the visual settings")
		return false
	}

	queueColorSignifier := cap.img.At(int(cap.xBoundForCalc*0.5901041666666667), int(cap.yBoundForCalc*0.4953703703703704))
	queueSelectingColorSignifier := cap.img.At(int(cap.xBoundForCalc*0.703125), int(cap.yBoundForCalc*0.4953703703703704))
	inQueueColorSignifier := cap.img.At(int(cap.xBoundForCalc*0.5729166666666667), 0)
	// fmt.Println(queueColorSignifier)
	// fmt.Println(inQueueColorSignifier)
	if ColorMatch(qpColor, queueColorSignifier, 15000) || ColorMatch(qpColor, queueSelectingColorSignifier, 15000) {
		c.app.Event.Emit("status-update", map[string]string{
			"statusIcon": "../public/assets/selection.svg",
			"statusText": "Selecting Role",
			"message":    "Select your role(s) and game mode and I'll lock them in once you queue.",
		})

		gameState.Selector.Queue = QP
		gameState.GameStatus = StatusSelection
		fmt.Println(gameState.GameStatus.String())
		fmt.Println(gameState.Selector.Queue.String())

		c.updateSelections(cap, gameState)
	} else if ColorMatch(compColor, queueColorSignifier, 15000) || ColorMatch(compColor, queueSelectingColorSignifier, 15000) {
		c.app.Event.Emit("status-update", map[string]string{
			"statusIcon": "../public/assets/selection.svg",
			"statusText": "Selecting Role",
			"message":    "Select your role(s) and game mode and I'll lock them in once you queue.",
		})

		gameState.Selector.Queue = Comp
		gameState.GameStatus = StatusSelection
		fmt.Println(gameState.GameStatus.String())
		fmt.Println(gameState.Selector.Queue.String())

		c.updateSelections(cap, gameState)
	} else {
		if ColorMatch(inQueueQPColor, inQueueColorSignifier, 15000) {
			c.app.Event.Emit("status-update", map[string]string{
				"statusIcon": "../public/assets/in-queue.svg",
				"statusText": "In Queue",
				"message":    "",
			})
			c.app.Event.Emit("queue-update", "Quick Play")
			c.app.Event.Emit("map-update", "Waiting for Match...")
			gameState.GameStatus = StatusInQueue
			gameState.Filters.GameMode = "0"
			c.estimateSelectionsOnEntry(gameState)
			fmt.Println(gameState.GameStatus.String())
		} else if ColorMatch(inQueueCompColor, inQueueColorSignifier, 15000) {
			c.app.Event.Emit("status-update", map[string]string{
				"statusIcon": "../public/assets/in-queue.svg",
				"statusText": "In Queue",
				"message":    "",
			})
			c.app.Event.Emit("queue-update", "Competitive")
			c.app.Event.Emit("map-update", "Waiting for Match...")
			gameState.GameStatus = StatusInQueue
			gameState.Filters.GameMode = "2"
			c.estimateSelectionsOnEntry(gameState)
			fmt.Println(gameState.GameStatus.String())
		} else {
			processedImg, _ := processImage(cap.img)
			text, err := analyze(processedImg)
			if err != nil {
				fmt.Printf("Error in Text Analysis: %e", err)
			}

			if strings.Contains(text, "VOTE") && strings.Contains(text, "MAP") {
				c.app.Event.Emit("status-update", map[string]string{
					"statusIcon": "../public/assets/map-voting.svg",
					"statusText": "Map Voting",
					"message":    "I see you've entered a match! I'll collect our final data needed after the map vote",
				})
				gameState.GameStatus = StatusMapVotingPhase
				fmt.Println(gameState.GameStatus.String())
			} else if strings.Contains(text, "RESULT") || strings.Contains(text, "SELECT") {
				gameState.GameStatus = StatusBanningPhase
				fmt.Println(gameState.GameStatus.String())
				c.app.Event.Emit("status-update", map[string]string{
					"statusIcon": "../public/assets/banning-phase.svg",
					"statusText": "Bans Phase",
					"message":    "Detecting map vote result...It may take a few seconds for me to find the name of the map, so don't worry if this takes a bit",
				})
				c.app.Event.Emit("map-update", "Detecting Map...")
			}
		}

	}

	return true
}

func (c *CaptureService) evaluate(counter int, gameState *GameState, done chan struct{}) {

	select {
	case <-c.stopChan:
		return
	default:
	}

	cap, err := capture()
	if err != nil {
		c.app.Event.Emit("error", err.Error()+"\n\nAn error occured while monitoring your screen, make sure Overwatch is running and set to 'Borderless Windowed' in the visual settings")
		select {
		case done <- struct{}{}:
			fmt.Println("Exiting...Capture Error")
			return // Exit evaluate immediately
		default:
			// Channel is full, someone else already sent
			return
		}
	}

	queueColorSignifier := cap.img.At(int(cap.xBoundForCalc*0.5901041666666667), int(cap.yBoundForCalc*0.4953703703703704))
	queueSelectingColorSignifier := cap.img.At(int(cap.xBoundForCalc*0.703125), int(cap.yBoundForCalc*0.4953703703703704))
	inQueueColorSignifier := cap.img.At(int(cap.xBoundForCalc*0.5729166666666667), 0)

	if ColorMatch(qpColor, queueColorSignifier, 15000) || ColorMatch(qpColor, queueSelectingColorSignifier, 15000) {
		if gameState.GameStatus != StatusSelection || gameState.Selector.Queue == Comp {
			c.app.Event.Emit("status-update", map[string]string{
				"statusIcon": "../public/assets/selection.svg",
				"statusText": "Selecting Role",
				"message":    "Select your role(s) and game mode and I'll lock them in once you queue.",
			})
			c.app.Event.Emit("queue-update", "Waiting...")
			c.app.Event.Emit("role-update", "Waiting...")

			gameState.Selector.Queue = QP
			gameState.GameStatus = StatusSelection
			fmt.Println(gameState.GameStatus.String())
			fmt.Println(gameState.Selector.Queue.String())
		}
		c.updateSelections(cap, gameState)
	} else if ColorMatch(compColor, queueColorSignifier, 15000) || ColorMatch(compColor, queueSelectingColorSignifier, 15000) {
		if gameState.GameStatus != StatusSelection || gameState.Selector.Queue == QP {
			c.app.Event.Emit("status-update", map[string]string{
				"statusIcon": "../public/assets/selection.svg",
				"statusText": "Selecting Role",
				"message":    "Select your role(s) and game mode and I'll will lock them in once you queue.",
			})
			c.app.Event.Emit("queue-update", "Waiting...")
			c.app.Event.Emit("role-update", "Waiting...")

			gameState.Selector.Queue = Comp
			gameState.GameStatus = StatusSelection
			fmt.Println(gameState.GameStatus.String())
			fmt.Println(gameState.Selector.Queue.String())
		}
		c.updateSelections(cap, gameState)
	} else {
		if gameState.GameStatus == StatusSelection || gameState.GameStatus == StatusIdle {
			if ColorMatch(inQueueQPColor, inQueueColorSignifier, 15000) {
				if gameState.GameStatus != StatusInQueue {
					c.app.Event.Emit("status-update", map[string]string{
						"statusIcon": "../public/assets/in-queue.svg",
						"statusText": "In Queue",
						"message":    defaultStatusMessage,
					})
					c.app.Event.Emit("queue-update", "Quick Play")
					c.app.Event.Emit("map-update", "Waiting for Match...")
					c.confirmSelections("0", gameState)
				}
			} else if ColorMatch(inQueueCompColor, inQueueColorSignifier, 15000) {
				if gameState.GameStatus != StatusInQueue {
					c.app.Event.Emit("status-update", map[string]string{
						"statusIcon": "../public/assets/in-queue.svg",
						"statusText": "In Queue",
						"message":    defaultStatusMessage,
					})
					c.app.Event.Emit("queue-update", "Competitive")
					c.app.Event.Emit("map-update", "Waiting for Match...")
					c.confirmSelections("2", gameState)
				}
			}
		} else if gameState.GameStatus == StatusInQueue {
			// enteringMatchAnalysisStart := time.Now()
			processedImg, _ := processImage(cap.img)
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
					c.app.Event.Emit("map-update", "Waiting for Map Vote...")
				}
				gameState.GameStatus = StatusMapVotingPhase
				fmt.Println(gameState.GameStatus.String())
			}
			// fmt.Printf("	Pre Match Analysis: %v\n", time.Since(enteringMatchAnalysisStart))
		} else if gameState.GameStatus == StatusMapVotingPhase {
			processedImg, _ := processImage(cap.img)
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
						"message":    "Detecting map vote result...It may take a few seconds for me to find the name of the map, so don't worry if this takes a bit",
					})
					c.app.Event.Emit("map-update", "Detecting Map...")
				}

				gameState.GameStatus = StatusBanningPhase
				postVoteText, err := analyze(processedImg)
				if err != nil {
					fmt.Printf("Error in Text Analysis: %e", err)
				}

				for i := range len(mapScan) {
					if strings.Contains(postVoteText, mapScan[i]) {
						mapDetection(gameState, mapScan[i], done)
						return
					}
				}

				exception := mapExceptions(postVoteText)

				if exception != "" {
					mapDetection(gameState, exception, done)
					return
				}
				fmt.Println(gameState.GameStatus.String())
			}
		} else if gameState.GameStatus == StatusBanningPhase {
			if gameState.Filters.GameMode == "0" {
				processedImg, _ := processImage(cap.img)
				text, err := analyze(processedImg)
				if err != nil {
					fmt.Printf("Error in Text Analysis: %e", err)
				}

				fmt.Println(text)

				for i := range len(mapScan) {
					if strings.Contains(text, mapScan[i]) {
						mapDetection(gameState, mapScan[i], done)
						return
					}
				}

				exception := mapExceptions(text)

				if exception != "" {
					mapDetection(gameState, exception, done)
					return
				}

			} else {
				// processedImg, _ := processImage(cap.img)
				// text, err := analyze(processedImg)
				// if err != nil {
				// 	fmt.Printf("Error in Text Analysis: %e", err)
				// }

				// fmt.Println(text)

				mapTopCap, mapBotCap, err := captureMap()
				if err != nil {
					fmt.Println(err)
					return
				}

				processedMapTopCap, _ := processImage(mapTopCap.img)
				mapTopText, err := analyze(processedMapTopCap)
				if err != nil {
					fmt.Printf("Error in Text Analysis: %e", err)
				}

				processedMapBotCap, _ := processImage(mapBotCap.img)
				mapBotText, err := analyze(processedMapBotCap)
				if err != nil {
					fmt.Printf("Error in Text Analysis: %e", err)
				}

				fmt.Println(mapTopText)
				fmt.Println(mapBotText)

				for i := range len(mapScan) {
					// if strings.Contains(text, mapScan[i]) {
					// 	gameState.Filters.Map = mapFormat[mapScan[i]]
					// 	select {
					// 	case done <- struct{}{}:
					// 		fmt.Printf("Map detected: %s\n", mapFormat[mapScan[i]])
					// 		return // Exit evaluate immediately
					// 	default:
					// 		// Channel is full, someone else already sent
					// 		return
					// 	}
					// }

					if strings.Contains(mapTopText, mapScan[i]) {
						mapDetection(gameState, mapScan[i], done)
						return
					}

					if strings.Contains(mapBotText, mapScan[i]) {
						mapDetection(gameState, mapScan[i], done)
						return
					}
				}

				exception := mapExceptions(mapTopText)

				if exception != "" {
					mapDetection(gameState, exception, done)
					return
				}

				exception = mapExceptions(mapBotText)

				if exception != "" {
					mapDetection(gameState, exception, done)
					return
				}
			}

		}
	}

	// if counter%10 == 0 {
	//     runtime.GC()
	// }
}

func (c *CaptureService) updateSelections(cap Capture, gameState *GameState) {
	err := roleRecognition(cap, gameState)
	if err != nil {
		fmt.Printf("Error in Role Recognition: %e", err)
	}
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
			gameState.Filters.Role = "Support"
			c.app.Event.Emit("message", "Oops! Looks like I couldn't detect your role. I will make you Support as a default right now, but once we have collected the other match data for you, you'll be able to pick what role you end up with in the filters")
			c.app.Event.Emit("role-update", "Support")
		}
	case StatusIdle:
		gameState.GameStatus = StatusInQueue
		gameState.Filters.GameMode = queueQueryParam
	}

	fmt.Println(gameState.Filters)
}

func (c *CaptureService) estimateSelectionsOnEntry(gameState *GameState) {
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
		gameState.Filters.Role = "Support"
		c.app.Event.Emit("message", "Oops! Looks like I couldn't detect your role. I will make you Support as a default right now, but once we have collected the other match data for you, you'll be able to pick what role you end up with in the filters")
		c.app.Event.Emit("role-update", "Support")
	}
}

func mapDetection(gameState *GameState, mapDetect string, done chan struct{}) {

	gameState.Filters.Map = mapFormat[mapDetect]
	select {
	case done <- struct{}{}:
		fmt.Printf("Map detected: %s\n", mapFormat[mapDetect])
		return // Exit evaluate immediately
	default:
		// Channel is full, someone else already sent
		return
	}

}

// Temporary Fix for Map Text Recognition Issues
func mapExceptions(text string) string {
	if strings.Contains(text, "AICHENWALDE") || strings.Contains(text, "HICHENWALDE") {
		return "EICHENWALDE"
	} else if strings.Contains(text, "ISPERANCA") || strings.Contains(text, "FSPERANCA") {
		return "ESPERANÇA"
	} else if strings.Contains(text, "LUIANG TOWER") {
		return "LIJIANG TOWER"
	} else {
		return ""
	}
}
