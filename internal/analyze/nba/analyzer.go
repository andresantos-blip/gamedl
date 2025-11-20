package nba

import (
	"encoding/json"
	"fmt"
	"gamedl/internal/common"
	"gamedl/lib/web/clients/sportsradar"
	"io"
	"os"
	"path/filepath"
)

type Analyzer struct {
	inputDir  string
	outputDir string
}

type ProcessResultNba struct {
	ID         string
	EventTypes map[string]int
	HasLane    bool
}

func NewAnalyzer(inputDir, outputDir string) *Analyzer {
	return &Analyzer{
		inputDir:  inputDir,
		outputDir: outputDir,
	}
}

func NewProcessResultNba() ProcessResultNba {
	return ProcessResultNba{
		EventTypes: make(map[string]int),
		HasLane:    false,
	}
}

func (a *Analyzer) processFileNba(path string) (ProcessResultNba, error) {
	result := NewProcessResultNba()
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return result, fmt.Errorf("could not open file %s: %w", path, err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return result, fmt.Errorf("could not read file %s: %w", path, err)
	}

	pbpData := &sportsradar.NbaGamePbp{}
	err = json.Unmarshal(data, pbpData)
	if err != nil {
		return result, fmt.Errorf("could not unmarshal game pbp: %w", err)
	}

	periods := pbpData.Periods
	for _, period := range periods {
		pbpEvents := period.Events
		for _, pbpEvent := range pbpEvents {
			eventType := pbpEvent.EventType
			result.EventTypes[eventType]++

			if eventType == "lane" || eventType == "doublelane" {
				result.HasLane = true
			}
		}
	}
	result.ID = pbpData.ID
	return result, nil
}

func (a *Analyzer) AnalyzeLaneViolations(years []int) error {
	var errs []error
	eventTypeCount := make(map[string]int)
	gamesWithLaneViolations := make(map[string]int) // gameID -> year

	for _, year := range years {
		path := common.GetYearGlobPattern(a.inputDir, "nba", year)
		matches, err := filepath.Glob(path)
		if err != nil {
			fmt.Printf("Error globbing files for year %d: %v\n", year, err)
			continue
		}

		fmt.Printf("year: %d, matches: %v\n", year, len(matches))
		for _, match := range matches {
			result, err := a.processFileNba(match)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			// Aggregate event type counts across all games
			for eventType, count := range result.EventTypes {
				eventTypeCount[eventType] += count
			}

			// Track games with lane violations
			if result.HasLane {
				gamesWithLaneViolations[result.ID] = year
			}
		}
	}

	// Write event type count report
	if err := a.writeJSONFile("event_type_count.json", eventTypeCount); err != nil {
		return fmt.Errorf("writing event_type_count: %w", err)
	}

	// Create lane violations games directory and copy games
	laneViolationsGamesDir := filepath.Join(a.outputDir, "lane_violations_games")
	if err := os.MkdirAll(laneViolationsGamesDir, 0755); err != nil {
		fmt.Printf("could not create lane_violations_games directory: %v\n", err)
	} else {
		for gameID, year := range gamesWithLaneViolations {
			gameFile := common.GetGameFilePath(a.inputDir, "nba", year, gameID)
			gameData, err := os.ReadFile(gameFile)
			if err != nil {
				fmt.Printf("could not read game file %s: %v\n", gameFile, err)
				continue
			}

			baseName := fmt.Sprintf("%s.json", gameID)
			err = os.WriteFile(filepath.Join(laneViolationsGamesDir, baseName), gameData, 0644)
			if err != nil {
				fmt.Printf("could not write game file: %v\n", err)
			} else {
				fmt.Printf("Copied game %s to lane_violations_games\n", gameID)
			}
		}
		fmt.Printf("Copied %d games with lane violations to %s\n", len(gamesWithLaneViolations), laneViolationsGamesDir)
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
