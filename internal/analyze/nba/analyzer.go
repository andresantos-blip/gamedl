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
	ID              string
	EventTypes      map[string]int
	TurnoverTypes   map[string]int
	HasLane         bool
	HasLaneViolationTurnover bool
	LaneViolations  []LaneViolationContext
}

type LaneViolationContext struct {
	Before []string `json:"before"`
	After  []string `json:"after"`
}

type GameLaneViolations struct {
	Year       int                    `json:"year"`
	GameID     string                 `json:"game_id"`
	Violations []LaneViolationContext `json:"violations"`
}

func NewAnalyzer(inputDir, outputDir string) *Analyzer {
	return &Analyzer{
		inputDir:  inputDir,
		outputDir: outputDir,
	}
}

func NewProcessResultNba() ProcessResultNba {
	return ProcessResultNba{
		EventTypes:                make(map[string]int),
		TurnoverTypes:             make(map[string]int),
		HasLane:                   false,
		HasLaneViolationTurnover:  false,
		LaneViolations:            make([]LaneViolationContext, 0),
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

	// Collect all events from all periods in order
	allEvents := make([]string, 0)
	periods := pbpData.Periods
	for _, period := range periods {
		pbpEvents := period.Events
		for _, pbpEvent := range pbpEvents {
			eventType := pbpEvent.EventType
			result.EventTypes[eventType]++
			allEvents = append(allEvents, eventType)

			if eventType == "lane" || eventType == "doublelane" {
				result.HasLane = true
			}

			// Track turnover_type counts
			if pbpEvent.TurnoverType != "" {
				result.TurnoverTypes[pbpEvent.TurnoverType]++
				if pbpEvent.TurnoverType == "Lane Violation" {
					result.HasLaneViolationTurnover = true
				}
			}
		}
	}

	// Find lane violations and their context
	for i, eventType := range allEvents {
		if eventType == "lane" || eventType == "doublelane" {
			// Get up to 5 events before
			beforeStart := max(0, i-5)
			beforeEvents := allEvents[beforeStart:i]
			
			// Get up to 5 events after
			afterEnd := min(len(allEvents), i+6)
			afterEvents := allEvents[i+1:afterEnd]

			result.LaneViolations = append(result.LaneViolations, LaneViolationContext{
				Before: beforeEvents,
				After:  afterEvents,
			})
		}
	}

	result.ID = pbpData.ID
	return result, nil
}

func (a *Analyzer) AnalyzeLaneViolations(years []int) error {
	var errs []error
	eventTypeCount := make(map[string]int)
	turnoverTypeCount := make(map[string]int)
	gamesWithLaneViolations := make(map[string]int) // gameID -> year
	gamesWithLaneViolationTurnovers := make(map[string]int) // gameID -> year
	gamesLaneViolationsContext := make([]GameLaneViolations, 0)

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

			// Aggregate turnover type counts across all games
			for turnoverType, count := range result.TurnoverTypes {
				turnoverTypeCount[turnoverType] += count
			}

			// Track games with lane violations (event_type)
			if result.HasLane {
				gamesWithLaneViolations[result.ID] = year
				// Store lane violations context for this game
				if len(result.LaneViolations) > 0 {
					gamesLaneViolationsContext = append(gamesLaneViolationsContext, GameLaneViolations{
						Year:       year,
						GameID:     result.ID,
						Violations: result.LaneViolations,
					})
				}
			}

			// Track games with "Lane Violation" turnover_type
			if result.HasLaneViolationTurnover {
				gamesWithLaneViolationTurnovers[result.ID] = year
			}
		}
	}

	// Write event type count report
	if err := a.writeJSONFile("event_type_count.json", eventTypeCount); err != nil {
		return fmt.Errorf("writing event_type_count: %w", err)
	}

	// Write turnover type count report
	if err := a.writeJSONFile("turnover_type_count.json", turnoverTypeCount); err != nil {
		return fmt.Errorf("writing turnover_type_count: %w", err)
	}

	// Write lane violations context document
	if err := a.writeJSONFile("lane_violations_context.json", gamesLaneViolationsContext); err != nil {
		return fmt.Errorf("writing lane_violations_context: %w", err)
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

	// Create lane violation turnover games directory and copy games
	laneViolationTurnoverGamesDir := filepath.Join(a.outputDir, "lane_violation_turnover_games")
	if err := os.MkdirAll(laneViolationTurnoverGamesDir, 0755); err != nil {
		fmt.Printf("could not create lane_violation_turnover_games directory: %v\n", err)
	} else {
		for gameID, year := range gamesWithLaneViolationTurnovers {
			gameFile := common.GetGameFilePath(a.inputDir, "nba", year, gameID)
			gameData, err := os.ReadFile(gameFile)
			if err != nil {
				fmt.Printf("could not read game file %s: %v\n", gameFile, err)
				continue
			}

			baseName := fmt.Sprintf("%s.json", gameID)
			err = os.WriteFile(filepath.Join(laneViolationTurnoverGamesDir, baseName), gameData, 0644)
			if err != nil {
				fmt.Printf("could not write game file: %v\n", err)
			} else {
				fmt.Printf("Copied game %s to lane_violation_turnover_games\n", gameID)
			}
		}
		fmt.Printf("Copied %d games with Lane Violation turnovers to %s\n", len(gamesWithLaneViolationTurnovers), laneViolationTurnoverGamesDir)
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
