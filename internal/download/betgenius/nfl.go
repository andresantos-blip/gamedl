package betgenius

import (
	"bytes"
	"encoding/json"
	"fmt"
	betgenius2 "gamedl/lib/web/clients/betgenius"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

const DefaultGamesDirectoryNfl = "nfl_games"

type GameProcessReport struct {
	Err  error
	Id   string
	Year int
}

func gamesPerYearNfl(client *betgenius2.Client, seasons *betgenius2.SeasonsReply) (map[int][]*betgenius2.Fixture, error) {
	years := seasons.SeasonsToYear()
	yearToGames := make(map[int][]*betgenius2.Fixture)

	for id, year := range years {
		schedule, err := client.GetNflGamesForSeason(id)
		if err != nil {
			return nil, fmt.Errorf("getting game schedule for year %v: %w", year, err)
		}
		yearToGames[year] = schedule.Embedded.Fixtures
	}

	return yearToGames, nil
}

func fetchAndSaveGameNfl(client *betgenius2.Client, gameID string, year int, outputDir string) error {
	gamePbpData, err := client.GetNflPbpRaw(gameID)
	if err != nil {
		return fmt.Errorf("fetching game pbp: %w", err)
	}

	gamesDir := filepath.Join(outputDir, DefaultGamesDirectoryNfl)
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

func DownloadNFL(seasons []int, concurrency int, outputDir string) error {
	client, err := createBetGeniusClient()
	if err != nil {
		return fmt.Errorf("failed to create BetGenius client: %w", err)
	}

	// Fetch seasons
	seasonsReply, err := client.GetNflSeasons("296")
	if err != nil {
		return fmt.Errorf("getting seasons: %w", err)
	}

	// Get games per year
	yearToGames, err := gamesPerYearNfl(client, seasonsReply)
	if err != nil {
		return fmt.Errorf("getting games: %w", err)
	}

	// Filter by requested seasons if specified
	if len(seasons) > 0 {
		filteredYearToGames := make(map[int][]*betgenius2.Fixture)
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
			gameStatus[game.StatusType]++
		}
		totalGames += gameStatus["scheduled"]
		fmt.Printf("Year: %v, Games: %v\n", year, len(games))
		for status, count := range gameStatus {
			fmt.Printf("  Status: %v, Count: %v\n", status, count)
		}
	}

	if totalGames == 0 {
		fmt.Println("No scheduled games found")
		return nil
	}

	tokenChannel := make(chan struct{}, concurrency)
	for i := 0; i < concurrency; i++ {
		tokenChannel <- struct{}{}
	}

	wg := sync.WaitGroup{}
	reportChannel := make(chan GameProcessReport, totalGames/concurrency+1)

	for year, games := range yearToGames {
		gamesDir := filepath.Join(outputDir, DefaultGamesDirectoryNfl)
		err := os.MkdirAll(filepath.Join(gamesDir, fmt.Sprintf("%d", year)), 0755)
		if err != nil {
			return fmt.Errorf("creating directory for year %d: %w", year, err)
		}

		for _, game := range games {
			if game.StatusType != "scheduled" {
				continue
			}
			wg.Add(1)
			go func(gameID int, gameYear int) {
				<-tokenChannel
				defer func() {
					tokenChannel <- struct{}{}
					wg.Done()
				}()

				report := GameProcessReport{
					Id:   strconv.Itoa(gameID),
					Year: gameYear,
				}

				fetchAndSaveError := fetchAndSaveGameNfl(client, report.Id, gameYear, outputDir)
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

		fmt.Printf("[%d] %s Processed game %s %d/%d (%.2f%%) games\n",
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
