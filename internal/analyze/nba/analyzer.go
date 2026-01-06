package nba

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gamedl/internal/common"
	"gamedl/lib/web/clients/sportsradar"
)

type Analyzer struct {
	inputDir  string
	outputDir string
}

type ProcessResultNba struct {
	ID                       string
	EventTypes               map[string]int
	TurnoverTypes            map[string]int
	HasLane                  bool
	HasLaneViolationTurnover bool
	LaneViolations           []LaneViolationContext
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

// PlayerStatsResult holds the result of processing a single game for player stats analysis
type PlayerStatsResult struct {
	GameID                      string
	EventTypesWithMissingPlayer map[string][]string
	HasMissingPlayerStats       bool
}

// GameMissingPlayerStats represents a game with events that have statistics without player info
type GameMissingPlayerStats struct {
	GameID     string   `json:"game_id"`
	EventTypes []string `json:"event_types"`
}

func NewAnalyzer(inputDir, outputDir string) *Analyzer {
	return &Analyzer{
		inputDir:  inputDir,
		outputDir: outputDir,
	}
}

func NewProcessResultNba() ProcessResultNba {
	return ProcessResultNba{
		EventTypes:               make(map[string]int),
		TurnoverTypes:            make(map[string]int),
		HasLane:                  false,
		HasLaneViolationTurnover: false,
		LaneViolations:           make([]LaneViolationContext, 0),
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
			afterEvents := allEvents[i+1 : afterEnd]

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
	gamesWithLaneViolations := make(map[string]int)         // gameID -> year
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
	if err := os.MkdirAll(laneViolationsGamesDir, 0o755); err != nil {
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
			err = os.WriteFile(filepath.Join(laneViolationsGamesDir, baseName), gameData, 0o644)
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
	if err := os.MkdirAll(laneViolationTurnoverGamesDir, 0o755); err != nil {
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
			err = os.WriteFile(filepath.Join(laneViolationTurnoverGamesDir, baseName), gameData, 0o644)
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
	if err := os.MkdirAll(a.outputDir, 0o755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}

	filePath := filepath.Join(a.outputDir, filename)
	if err := os.WriteFile(filePath, jsonData, 0o644); err != nil {
		return fmt.Errorf("writing file %s: %w", filePath, err)
	}

	fmt.Printf("Written: %s\n", filePath)
	return nil
}

// processFileForPlayerStats processes a single PBP file and returns statistics about events
// that have statistics without player information. The gameID is the path relative to inputDir.
func (a *Analyzer) processFileForPlayerStats(path string) (PlayerStatsResult, error) {
	// Derive game ID from relative path to inputDir
	gameID, err := filepath.Rel(a.inputDir, path)
	if err != nil {
		// Fallback to full path if relative path fails
		gameID = path
	}

	result := PlayerStatsResult{
		GameID:                      gameID,
		EventTypesWithMissingPlayer: make(map[string][]string),
		HasMissingPlayerStats:       false,
	}

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

	// Iterate through all periods and events
	for _, period := range pbpData.Periods {
		for _, event := range period.Events {
			// Check if this event has statistics
			if len(event.Statistics) == 0 {
				continue
			}

			// Check if any statistic is missing player info
			for _, stat := range event.Statistics {
				if stat.Player == nil {
					result.EventTypesWithMissingPlayer[event.EventType] = append(result.EventTypesWithMissingPlayer[event.EventType], event.ID)
					result.HasMissingPlayerStats = true
				}
			}
		}
	}

	return result, nil
}

// AnalyzePlayerStats analyzes PBP data to find events with statistics that are missing player information.
// It recursively searches for all JSON files in the input directory.
func (a *Analyzer) AnalyzePlayerStats() error {
	var errs []error

	// Aggregate counts of event types with missing player stats
	eventTypeCount := make(map[string]int)
	eventTypeCountUnique := make(map[string]int)

	// Track games with missing player stats and their affected event types
	gamesWithMissingPlayerStats := make([]GameMissingPlayerStats, 0)

	// filepath.Glob doesn't support ** for recursive matching, so we use filepath.WalkDir instead
	var matches []string
	err := filepath.WalkDir(a.inputDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".json" {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking directory: %w", err)
	}

	fmt.Printf("Found %d JSON files in %s\n", len(matches), a.inputDir)

	uniqueEventIds := make(map[string]struct{})

	for _, match := range matches {
		result, err := a.processFileForPlayerStats(match)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		// Aggregate event type counts
		for eventType, ids := range result.EventTypesWithMissingPlayer {
			eventTypeCount[eventType] += len(ids)
			for _, id := range ids {
				if _, ok := uniqueEventIds[id]; !ok {
					eventTypeCountUnique[eventType] += 1
					uniqueEventIds[id] = struct{}{}
				}
			}

		}

		// Track games with missing player stats
		if result.HasMissingPlayerStats {
			// Collect unique event types for this game
			eventTypes := make([]string, 0, len(result.EventTypesWithMissingPlayer))
			for eventType := range result.EventTypesWithMissingPlayer {
				eventTypes = append(eventTypes, eventType)
			}

			gamesWithMissingPlayerStats = append(gamesWithMissingPlayerStats, GameMissingPlayerStats{
				GameID:     result.GameID,
				EventTypes: eventTypes,
			})
		}
	}

	// Write event type count report
	if err := a.writeJSONFile("event_types_without_player_stats.json", eventTypeCount); err != nil {
		return fmt.Errorf("writing event_types_without_player_stats: %w", err)
	}

	if err := a.writeJSONFile("event_types_without_player_stats_unique.json", eventTypeCountUnique); err != nil {
		return fmt.Errorf("writing event_types_without_player_stats_unique: %w", err)
	}

	// Write games with missing player stats report
	if err := a.writeJSONFile("games_with_missing_player_stats.json", gamesWithMissingPlayerStats); err != nil {
		return fmt.Errorf("writing games_with_missing_player_stats: %w", err)
	}

	if len(errs) > 0 {
		fmt.Printf("Encountered %d errors during processing:\n", len(errs))
		for _, err := range errs {
			fmt.Printf("  %v\n", err)
		}
	}

	fmt.Printf("\nAnalysis complete:\n")
	fmt.Printf("  - Total event types with missing player stats: %d\n", len(eventTypeCount))
	fmt.Printf("  - Total games with missing player stats: %d\n", len(gamesWithMissingPlayerStats))

	return nil
}
