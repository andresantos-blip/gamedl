package ncaab

import (
	"encoding/json"
	"fmt"
	"gamedl/lib/sportsradar"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

type Analyzer struct {
	inputDir  string
	outputDir string
}

type ProcessResultNcaab struct {
	ID          string
	EventTypes  map[string]int
	BeforeEvent map[string][][]string
}

type GameReview struct {
	Year   int        `json:"year"`
	ID     string     `json:"id"`
	Before [][]string `json:"before"`
}

var NcaabReviewTypes = []string{
	"challengereview",
	"challengetimeout",
	"requestreview",
	"review",
}

func NewAnalyzer(inputDir, outputDir string) *Analyzer {
	return &Analyzer{
		inputDir:  inputDir,
		outputDir: outputDir,
	}
}

func NewProcessResultNcaab() ProcessResultNcaab {
	return ProcessResultNcaab{
		EventTypes:  make(map[string]int),
		BeforeEvent: make(map[string][][]string),
	}
}

func (a *Analyzer) processFileNcaab(path string) (ProcessResultNcaab, error) {
	result := NewProcessResultNcaab()
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return result, fmt.Errorf("could not open file %s: %w", path, err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return result, fmt.Errorf("could not read file %s: %w", path, err)
	}

	pbpData := &sportsradar.NcaabGamePbp{}
	err = json.Unmarshal(data, pbpData)
	if err != nil {
		return result, fmt.Errorf("could not unmarshal game pbp: %w", err)
	}

	periods := pbpData.Periods
	for _, period := range periods {
		pbpEvents := period.Events
		for i, pbpEvent := range pbpEvents {
			eventType := pbpEvent.EventType
			result.EventTypes[eventType]++

			if slices.Contains(NcaabReviewTypes, eventType) {
				beforeEvent := make([]string, 0, i)
				for j := 0; j < i; j++ {
					beforeEvent = append(beforeEvent, pbpEvents[j].EventType)
				}
				result.BeforeEvent[eventType] = append(result.BeforeEvent[eventType], beforeEvent)
			}
		}
	}
	result.ID = pbpData.ID
	return result, nil
}

func (a *Analyzer) AnalyzeReviewTypes(years []int) error {
	defaultYears := []int{2012, 2013, 2014, 2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024}
	if len(years) == 0 {
		years = defaultYears
	}

	var errs []error
	eventsToGames := make(map[string][]*GameReview)
	eventTypeCount := make(map[string]int)

	for _, year := range years {
		path := filepath.Join(a.inputDir, "ncaab_games", strconv.Itoa(year), "*.json")
		matches, err := filepath.Glob(path)
		if err != nil {
			fmt.Printf("Error globbing files for year %d: %v\n", year, err)
			continue
		}

		fmt.Printf("year: %d, matches: %v\n", year, len(matches))
		for _, match := range matches {
			result, err := a.processFileNcaab(match)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			for eventType, count := range result.EventTypes {
				eventTypeCount[eventType] += count
				if slices.Contains(NcaabReviewTypes, eventType) {
					eventsToGames[eventType] = append(
						eventsToGames[eventType],
						&GameReview{Year: year, ID: result.ID, Before: result.BeforeEvent[eventType]},
					)
				}
			}
		}
	}

	// Trim "before" events to last 5
	for t, games := range eventsToGames {
		for i := range games {
			for j := range eventsToGames[t][i].Before {
				start := max(len(eventsToGames[t][i].Before[j])-5, 0)
				eventsToGames[t][i].Before[j] = eventsToGames[t][i].Before[j][start:]
			}
		}
	}

	// Write results
	if err := a.writeJSONFile("review_events_to_games.json", eventsToGames); err != nil {
		return fmt.Errorf("writing review_events_to_games: %w", err)
	}

	if err := a.writeJSONFile("event_type_count.json", eventTypeCount); err != nil {
		return fmt.Errorf("writing event_type_count: %w", err)
	}

	// Create review games directory and copy sample games
	reviewGamesDir := filepath.Join(a.outputDir, "review_games_ncaab")
	if err := os.MkdirAll(reviewGamesDir, 0755); err != nil {
		fmt.Printf("could not create review_games directory: %v\n", err)
	} else {
		for eventType, games := range eventsToGames {
			if len(games) == 0 {
				continue
			}

			for _, game := range games {
				gameFile := filepath.Join(a.inputDir, "ncaab_games", strconv.Itoa(game.Year), fmt.Sprintf("%s.json", game.ID))
				gameData, err := os.ReadFile(gameFile)
				if err != nil {
					fmt.Printf("could not read game file: %v\n", err)
					continue
				}

				cleanEventType := eventType
				if cleanEventType == "" {
					cleanEventType = "no_type"
				} else {
					cleanEventType = strings.ReplaceAll(cleanEventType, " ", "_")
				}

				baseName := fmt.Sprintf("%s-%s.json", cleanEventType, game.ID)
				err = os.WriteFile(filepath.Join(reviewGamesDir, baseName), gameData, 0644)
				if err != nil {
					fmt.Printf("could not write game file: %v\n", err)
				}
			}
		}
	}

	if len(errs) > 0 {
		fmt.Printf("Encountered %d errors during processing:\n", len(errs))
		for _, err := range errs {
			fmt.Printf("  %v\n", err)
		}
	}

	return nil
}

func (a *Analyzer) writeJSONFile(filename string, data interface{}) error {
	// Ensure output directory exists
	if err := os.MkdirAll(a.outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}

	filePath := filepath.Join(a.outputDir, filename)
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("writing file %s: %w", filePath, err)
	}

	fmt.Printf("Written: %s\n", filePath)
	return nil
}
