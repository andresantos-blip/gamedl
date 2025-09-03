package sportradar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gamedl/lib/sportsradar"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

const DefaultGamesDirectoryNcaaf = "ncaaf_games"

func gamesPerYearNcaaf(client *sportsradar.Client, seasons *sportsradar.NcaafSeasonsInfo) (map[int][]*sportsradar.NcaafGame, error) {
	years := seasons.Years()
	yearToGames := make(map[int][]*sportsradar.NcaafGame)

	for _, year := range years {
		schedule, err := client.GetNcaafSeasonSchedule(year)
		if err != nil {
			return nil, fmt.Errorf("getting game schedule for year %v: %w", year, err)
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

func fetchAndSaveGameNcaaf(client *sportsradar.Client, gameID string, year int, outputDir string) error {
	gamePbpData, err := client.GetNcaafPbpOfGameRaw(gameID)
	if err != nil {
		return fmt.Errorf("fetching game pbp: %w", err)
	}

	gamesDir := filepath.Join(outputDir, DefaultGamesDirectoryNcaaf)
	pathtoFile := filepath.Join(gamesDir, strconv.Itoa(year), fmt.Sprintf("%s.json", gameID))

	bytesBuffer := bytes.NewBuffer([]byte{})
	err = json.Indent(bytesBuffer, gamePbpData, "", "  ")
	if err != nil {
		return fmt.Errorf("indenting game pbp: %w", err)
	}

	defer func() {
		bytesBuffer.Reset()
		bytesBuffer = nil
	}()

	err = os.WriteFile(pathtoFile, bytesBuffer.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("saving game pbp: %w", err)
	}

	return nil
}

func DownloadNCAF(seasons []int, concurrency int, outputDir string) error {
	client, err := createSportRadarClientWithNCAF()
	if err != nil {
		return fmt.Errorf("failed to create SportRadar client: %w", err)
	}

	// Fetch seasons
	seasonsInfo, err := client.GetNcaafSeasons()
	if err != nil {
		return fmt.Errorf("getting seasons: %w", err)
	}

	// Get games per year
	yearToGames, err := gamesPerYearNcaaf(client, seasonsInfo)
	if err != nil {
		return fmt.Errorf("getting games: %w", err)
	}

	// Filter by requested seasons if specified
	if len(seasons) > 0 {
		filteredYearToGames := make(map[int][]*sportsradar.NcaafGame)
		for _, season := range seasons {
			if games, exists := yearToGames[season]; exists {
				filteredYearToGames[season] = games
			}
		}
		yearToGames = filteredYearToGames
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

	if totalGames == 0 {
		fmt.Println("No closed games found")
		return nil
	}

	tokenChannel := make(chan struct{}, concurrency)
	for i := 0; i < concurrency; i++ {
		tokenChannel <- struct{}{}
	}

	wg := sync.WaitGroup{}
	reportChannel := make(chan GameProcessReport, totalGames/concurrency+1)

	for year, games := range yearToGames {
		gamesDir := filepath.Join(outputDir, DefaultGamesDirectoryNcaaf)
		err := os.MkdirAll(filepath.Join(gamesDir, fmt.Sprintf("%d", year)), 0755)
		if err != nil {
			return fmt.Errorf("creating directory for year %d: %w", year, err)
		}

		for _, game := range games {
			if game.Status != "closed" {
				continue
			}
			wg.Add(1)
			go func(gameID string, gameYear int) {
				<-tokenChannel
				defer func() {
					tokenChannel <- struct{}{}
					wg.Done()
				}()

				report := GameProcessReport{
					Id:   gameID,
					Year: gameYear,
				}

				fetchAndSaveError := fetchAndSaveGameNcaaf(client, gameID, gameYear, outputDir)
				if fetchAndSaveError != nil {
					report.Err = fetchAndSaveError
				}
				reportChannel <- report
			}(game.ID, year)
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

		fmt.Printf("[%d] %s Downloaded game %s | Progress: %d/%d (%.2f%%) games\n",
			report.Year, status, report.Id, processed, totalGames, (float64(processed)/float64(totalGames))*100.0)
	}

	if len(reportErrors) > 0 {
		fmt.Printf("Errors:\n")
		for _, report := range reportErrors {
			fmt.Printf("  %s: %v\n", report.Id, report.Err)
		}
	}

	return nil
}
