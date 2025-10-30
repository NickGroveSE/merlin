package main

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

	heroes, _ := c.overwatchService.Scrape(captureFilters)

	captureHandler()

	return heroes, captureFilters, nil
}
