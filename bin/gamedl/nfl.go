package main

import (
	"fmt"
	"gamedl/lib/betgenius"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

const GamesDirectoryNfl = "nfl_games"

func GamesPerYearNfl(client *betgenius.Client, seasons *betgenius.SeasonsReply) (map[int][]*betgenius.Fixture, error) {
	years := seasons.SeasonsToYear()
	yearToGames := make(map[int][]*betgenius.Fixture)

	for id, year := range years {
		schedule, err := client.GetNflGamesForSeason(id)
		if err != nil {
			return nil, fmt.Errorf("getting game schedule for year %v: %w\n", year, err)
		}

		yearToGames[year] = schedule.Embedded.Fixtures

	}

	return yearToGames, nil
}

func NflGames() {
	fixtureKey := os.Getenv("BG_FIXTURE_KEY")
	if fixtureKey == "" {
		fmt.Printf("BG_FIXTURE_KEY environment variable not set\n")
		os.Exit(1)
	}

	fixtureUsername := os.Getenv("BG_FIXTURE_USER")
	if fixtureUsername == "" {
		fmt.Printf("BG_FIXTURE_USER environment variable not set\n")
		os.Exit(1)
	}

	fixturePassword := os.Getenv("BG_FIXTURE_PASSWORD")
	if fixturePassword == "" {
		fmt.Printf("BG_FIXTURE_PASSWORD environment variable not set\n")
		os.Exit(1)
	}

	statsKey := os.Getenv("BG_STATS_KEY")
	if statsKey == "" {
		fmt.Printf("BG_STATS_KEY environment variable not set\n")
		os.Exit(1)
	}

	statsUsername := os.Getenv("BG_STATS_USER")
	if statsUsername == "" {
		fmt.Printf("BG_STATS_USER environment variable not set\n")
		os.Exit(1)
	}

	statsPassword := os.Getenv("BG_STATS_PASSWORD")
	if statsPassword == "" {
		fmt.Printf("BG_STATS_PASSWORD environment variable not set\n")
		os.Exit(1)
	}

	client := betgenius.NewClient(
		betgenius.WithStatsKey(statsKey),
		betgenius.WithFixtureUsername(fixtureUsername),
		betgenius.WithFixturePassword(fixturePassword),
		betgenius.WithFixtureKey(fixtureKey),
		betgenius.WithStatsUsername(statsUsername),
		betgenius.WithStatsPassword(statsPassword),
	)

	// Fetch seasons
	seasons, err := client.GetNflSeasons("296")
	if err != nil {
		fmt.Printf("Error getting seasons: %v\n", err)
		os.Exit(1)
	}

	// Get games per year
	yearToGames, err := GamesPerYearNfl(client, seasons)
	if err != nil {
		fmt.Printf("Error getting games: %v\n", err)
		os.Exit(1)
	}

	totalGames := 0
	for year, games := range yearToGames {
		gameStatus := make(map[string]int)
		for _, game := range games {
			gameStatus[game.StatusType]++
		}
		totalGames += gameStatus["scheduled"]
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
		err := os.MkdirAll(filepath.Join(GamesDirectoryNfl, fmt.Sprintf("%d", year)), 0755)
		if err != nil {
			fmt.Printf("Error creating directory for year %d: %v\n", year, err)
			os.Exit(1)
		}
		for _, game := range games {
			if game.StatusType != "scheduled" {
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
					Id:   strconv.Itoa(game.ID),
					Year: year,
				}
				fetchAndSaveError := FetchAndSaveGameNfl(client, report.Id, year)
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

func FetchAndSaveGameNfl(client *betgenius.Client, gameID string, year int) error {
	// Fetch game
	gamePbpData, err := client.GetNflPbpRaw(gameID)
	if err != nil {
		return fmt.Errorf("fetching game pbp: %w", err)
	}

	// Save game
	pathtoFile := filepath.Join(GamesDirectoryNfl, strconv.Itoa(year), fmt.Sprintf("%s.json", gameID))
	err = os.WriteFile(pathtoFile, gamePbpData, 0644)
	if err != nil {
		return fmt.Errorf("saving game pbp: %w", err)
	}

	return nil
}
