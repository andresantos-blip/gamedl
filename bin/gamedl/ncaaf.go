package main

import (
	"fmt"
	"gamedl/lib/sportsradar"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

const GamesDirectoryNccaf = "ncaaf_games"

func GamesPerYearNcaaf(client *sportsradar.Client, seasons *sportsradar.NcaafSeasonsInfo) (map[int][]*sportsradar.NcaafGame, error) {
	years := seasons.Years()
	yearToGames := make(map[int][]*sportsradar.NcaafGame)

	for _, year := range years {
		schedule, err := client.GetNcaafSeasonSchedule(year)
		if err != nil {
			return nil, fmt.Errorf("getting game schedule for year %v: %w\n", year, err)
		}

		yearToGames[year] = make([]*sportsradar.NcaafGame, 0, 1024)
		for _, week := range schedule.Weeks {
			for _, game := range week.Games {
				yearToGames[year] = append(yearToGames[year], game)
			}
		}

	}

	return yearToGames, nil
}

func NcaafGames() {
	apiKey := os.Getenv("SPORTRADAR_NCAAF_KEY")
	if apiKey == "" {
		fmt.Printf("SPORTRADAR_NCAAF_KEY environment variable not set\n")
		os.Exit(1)
	}

	client := sportsradar.NewClient(sportsradar.WithNcaafAPIKey(apiKey))

	// Fetch seasons
	seasons, err := client.GetNcaafSeasons()
	if err != nil {
		fmt.Printf("Error getting seasons: %v\n", err)
		os.Exit(1)
	}

	// Get games per year
	yearToGames, err := GamesPerYearNcaaf(client, seasons)
	if err != nil {
		fmt.Printf("Error getting games: %v\n", err)
		os.Exit(1)
	}

	totalGames := 0
	for year, games := range yearToGames {
		gameStatus := make(map[string]int)
		for _, game := range games {
			gameStatus[game.Status]++
		}
		totalGames += gameStatus["closed"]
		fmt.Printf("Year: %v, Games: %v\n", year, len(games))
		for status, count := range gameStatus {
			fmt.Printf("  Status: %v, Count: %v\n", status, count)
		}
	}

	tokenChannel := make(chan struct{}, MaxConcurrency)
	for i := 0; i < MaxConcurrency; i++ {
		tokenChannel <- struct{}{}
	}

	wg := sync.WaitGroup{}
	reportChannel := make(chan GameProcessReport, totalGames/MaxConcurrency)

	for year, games := range yearToGames {
		err := os.MkdirAll(filepath.Join(GamesDirectoryNccaf, fmt.Sprintf("%d", year)), 0755)
		if err != nil {
			fmt.Printf("Error creating directory for year %d: %v\n", year, err)
			os.Exit(1)
		}
		for _, game := range games {
			if game.Status != "closed" {
				continue
			}
			wg.Add(1)
			go func() {
				<-tokenChannel
				defer func() {
					tokenChannel <- struct{}{}
					wg.Done()
				}()
				report := GameProcessReport{
					Id:   game.ID,
					Year: year,
				}
				fetchAndSaveError := FetchAndSaveGameNcaaf(client, game.ID, year)
				if fetchAndSaveError != nil {
					report.Err = fetchAndSaveError
				}
				reportChannel <- report
			}()
		}
	}

	go func() {
		wg.Wait()
		close(reportChannel)
	}()

	processed := 0
	var reportErrors []GameProcessReport

	for report := range reportChannel {
		processed++
		status := "✅"
		if report.Err != nil {
			reportErrors = append(reportErrors, report)
			status = "❌"
		}

		fmt.Printf("[%d] %s Processed game %s %d/%d (%.2f%%) games\n",
			report.Year, status, report.Id, processed, totalGames, (float64(processed)/float64(totalGames))*100.0)
	}

	if len(reportErrors) > 0 {
		fmt.Printf("Errors:\n")
		for _, report := range reportErrors {
			fmt.Printf("  %s: %v\n", report.Id, report.Err)
		}
	}

}

func FetchAndSaveGameNcaaf(client *sportsradar.Client, gameID string, year int) error {
	// Fetch game
	gamePbpData, err := client.GetNcaafPbpOfGameRaw(gameID)
	if err != nil {
		return fmt.Errorf("fetching game pbp: %w", err)
	}

	// Save game
	pathtoFile := filepath.Join(GamesDirectoryNccaf, strconv.Itoa(year), fmt.Sprintf("%s.json", gameID))
	err = os.WriteFile(pathtoFile, gamePbpData, 0644)
	if err != nil {
		return fmt.Errorf("saving game pbp: %w", err)
	}

	return nil
}
