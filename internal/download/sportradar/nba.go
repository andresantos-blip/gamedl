package sportradar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gamedl/internal/common"
	sportsradar2 "gamedl/lib/web/clients/sportsradar"
	"os"
	"sync"
)

func gamesPerYearNBA(client *sportsradar2.Client, seasons *sportsradar2.NBASeasonsInfo) (map[int][]*sportsradar2.NBAGame, error) {
	years := seasons.Years()
	yearToGames := make(map[int][]*sportsradar2.NBAGame)

	for _, year := range years {
		schedule, err := client.GetNbaSeasonSchedule(year)
		if err != nil {
			return nil, fmt.Errorf("getting game schedule for year %v: %w", year, err)
		}

		yearToGames[year] = make([]*sportsradar2.NBAGame, 0, 1024)
		for _, game := range schedule.Games {
			yearToGames[year] = append(yearToGames[year], game)
		}

	}

	return yearToGames, nil
}

func fetchAndSaveGameNBA(client *sportsradar2.Client, gameID string, year int, outputDir string) error {
	gamePbpData, err := client.GetNbaPbpOfGameRaw(gameID)
	if err != nil {
		return fmt.Errorf("fetching game pbp: %w", err)
	}

	pathtoFile := common.GetGameFilePath(outputDir, "NBA", year, gameID)

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

func DownloadNBA(seasons []int, concurrency int, outputDir string) error {
	client, err := createSportRadarClientWithNba()
	if err != nil {
		return fmt.Errorf("failed to create SportRadar client: %w", err)
	}

	fmt.Println("Getting seasons data...")
	// Fetch seasons
	seasonsInfo, err := client.GetNbaSeasons()
	if err != nil {
		return fmt.Errorf("getting seasons: %w", err)
	}

	if len(seasons) > 0 {
		seasonsInfo.FilterYears(seasons)
	}
	seasonsInfo.FilterSeasonType("REG")

	fmt.Printf("Getting game ids for seasons %v...\n", seasonsInfo.Years())

	// Get games per year
	yearToGames, err := gamesPerYearNBA(client, seasonsInfo)
	if err != nil {
		return fmt.Errorf("getting games: %w", err)
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
		err := common.CreateYearDirectory(outputDir, "NBA", year)
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

				fetchAndSaveError := fetchAndSaveGameNBA(client, gameID, gameYear, outputDir)
				if fetchAndSaveError != nil {
					report.Err = fetchAndSaveError
				}
				reportChannel <- report
			}(game.Id, year)
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
