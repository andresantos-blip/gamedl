package main

import (
	"fmt"
	"gamedl/lib/sportsradar"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

const GamesDirectoryNccab = "ncaab_games"

func GamesPerYearNcaab(client *sportsradar.Client, seasons *sportsradar.NcaabSeasonsInfo) (map[int][]*sportsradar.NcaabGame, error) {
	years := seasons.Years()
	yearToGames := make(map[int][]*sportsradar.NcaabGame)

	for _, year := range years {
		schedule, err := client.GetNcaabSeasonSchedule(year)
		if err != nil {
			return nil, fmt.Errorf("getting game schedule for year %v: %w\n", year, err)
		}

		yearToGames[year] = make([]*sportsradar.NcaabGame, 0, 1024)
		for _, game := range schedule.Games {
			yearToGames[year] = append(yearToGames[year], game)
		}

	}

	return yearToGames, nil
}

func NcaabGames() {
	apiKey := os.Getenv("SPORTRADAR_NCAAB_KEY")
	if apiKey == "" {
		fmt.Printf("SPORTRADAR_NCAAB_KEY environment variable not set\n")
		os.Exit(1)
	}

	client := sportsradar.NewClient(sportsradar.WithNcaabAPIKey(apiKey))

	// Fetch seasons
	seasons, err := client.GetNcaabSeasons()
	if err != nil {
		fmt.Printf("Error getting seasons: %v\n", err)
		os.Exit(1)
	}

	// Get games per year
	yearToGames, err := GamesPerYearNcaab(client, seasons)
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
		err := os.MkdirAll(filepath.Join(GamesDirectoryNccab, fmt.Sprintf("%d", year)), 0755)
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
				fetchAndSaveError := FetchAndSaveGameNcaab(client, game.ID, year)
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
			fmt.Printf("Error: %v\n", report.Err)
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

func FetchAndSaveGameNcaab(client *sportsradar.Client, gameID string, year int) error {
	// Fetch game
	gamePbpData, err := client.GetNcaabPbpOfGameRaw(gameID)
	if err != nil {
		return fmt.Errorf("fetching game pbp: %w", err)
	}

	// Save game
	pathtoFile := filepath.Join(GamesDirectoryNccab, strconv.Itoa(year), fmt.Sprintf("%s.json", gameID))
	err = os.WriteFile(pathtoFile, gamePbpData, 0644)
	if err != nil {
		return fmt.Errorf("saving game pbp: %w", err)
	}

	return nil
}
