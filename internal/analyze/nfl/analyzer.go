package nfl

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"

	"gamedl/internal/common"
	"gamedl/lib/web/clients/betgenius"
)

type Analyzer struct {
	inputDir  string
	outputDir string
}

type ProcessResultNfl struct {
	ID                      string
	ActionTypes             map[string]int
	ActionSubTypes          map[string]int
	BeforeAction            map[string][][]string
	RecoveriesInConversions []string
}

type GameReview struct {
	Year   int        `json:"year"`
	ID     string     `json:"id"`
	Before [][]string `json:"before"`
}

func NewAnalyzer(inputDir, outputDir string) *Analyzer {
	return &Analyzer{
		inputDir:  inputDir,
		outputDir: outputDir,
	}
}

func NewProcessResultNfl() ProcessResultNfl {
	return ProcessResultNfl{
		ActionTypes:             make(map[string]int),
		ActionSubTypes:          make(map[string]int),
		BeforeAction:            make(map[string][][]string),
		RecoveriesInConversions: make([]string, 0),
	}
}

func (a *Analyzer) processFileNfl(path string) (ProcessResultNfl, error) {
	result := NewProcessResultNfl()
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return result, fmt.Errorf("could not open file %s: %w", path, err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return result, fmt.Errorf("could not read file %s: %w", path, err)
	}

	pbpData := &betgenius.GamePbp{}
	err = json.Unmarshal(data, pbpData)
	if err != nil {
		return result, fmt.Errorf("could not unmarshal game pbp: %w", err)
	}

	drives := slices.Concat(pbpData.FirstHalf.Drives, pbpData.SecondHalf.Drives)
	for _, otPeriod := range pbpData.OvertimePeriods {
		drives = slices.Concat(drives, otPeriod.Drives)
	}

	for _, drive := range drives {
		for _, conversionPlay := range drive.ConversionPlays {
			cmi := slices.IndexFunc(conversionPlay.Actions, func(action betgenius.ConversionPlayAction) bool {
				return action.Type == "ConversionMade"
			})
			cmr := slices.IndexFunc(conversionPlay.Actions, func(action betgenius.ConversionPlayAction) bool {
				return action.Type == "Recovery"
			})

			if cmi != -1 && cmr != -1 && cmi > cmr {
				result.RecoveriesInConversions = append(result.RecoveriesInConversions, conversionPlay.ID)
			}

		}
		for _, play := range drive.Plays {
			for i, action := range play.Actions {
				result.ActionTypes[action.Type]++
				if action.SubType != nil {
					result.ActionSubTypes[*action.SubType]++
				}
				beforeAction := make([]string, 0, i)
				for j := 0; j < i; j++ {
					beforeAction = append(beforeAction, play.Actions[j].Type)
				}
				result.BeforeAction[action.Type] = append(result.BeforeAction[action.Type], beforeAction)
			}
		}
	}

	result.ID = pbpData.FixtureID
	return result, nil
}

func (a *Analyzer) AnalyzeActionTypes(years []int) error {
	var errs []error
	actionsToGames := make(map[string][]*GameReview)
	subActionsToGames := make(map[string][]*GameReview)
	actionTypeCount := make(map[string]int)
	subActionTypeCount := make(map[string]int)

	for _, year := range years {
		path := common.GetYearGlobPattern(a.inputDir, "nfl", year)
		matches, err := filepath.Glob(path)
		if err != nil {
			fmt.Printf("Error globbing files for year %d: %v\n", year, err)
			continue
		}

		fmt.Printf("year: %d, matches: %v\n", year, len(matches))
		for _, match := range matches {
			result, err := a.processFileNfl(match)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			for actionType, count := range result.ActionTypes {
				actionTypeCount[actionType] += count
				actionsToGames[actionType] = append(
					actionsToGames[actionType],
					&GameReview{Year: year, ID: result.ID, Before: result.BeforeAction[actionType]},
				)
			}

			for subActionType, count := range result.ActionSubTypes {
				subActionTypeCount[subActionType] += count
				subActionsToGames[subActionType] = append(
					subActionsToGames[subActionType],
					&GameReview{Year: year, ID: result.ID, Before: result.BeforeAction[subActionType]},
				)
			}
		}
	}

	// Trim "before" actions to last 10
	for t, games := range actionsToGames {
		for i := range games {
			for j := range actionsToGames[t][i].Before {
				start := max(len(actionsToGames[t][i].Before[j])-10, 0)
				actionsToGames[t][i].Before[j] = actionsToGames[t][i].Before[j][start:]
			}
		}
	}

	for t, games := range subActionsToGames {
		for i := range games {
			for j := range subActionsToGames[t][i].Before {
				start := max(len(subActionsToGames[t][i].Before[j])-10, 0)
				subActionsToGames[t][i].Before[j] = subActionsToGames[t][i].Before[j][start:]
			}
		}
	}

	// Write results
	if err := a.writeJSONFile("actions_to_games.json", actionsToGames); err != nil {
		return fmt.Errorf("writing actions_to_games: %w", err)
	}

	if err := a.writeJSONFile("sub_actions_to_games.json", subActionsToGames); err != nil {
		return fmt.Errorf("writing sub_actions_to_games: %w", err)
	}

	if err := a.writeJSONFile("action_type_count.json", actionTypeCount); err != nil {
		return fmt.Errorf("writing action_type_count: %w", err)
	}

	if err := a.writeJSONFile("sub_action_type_count.json", subActionTypeCount); err != nil {
		return fmt.Errorf("writing sub_action_type_count: %w", err)
	}

	if len(errs) > 0 {
		fmt.Printf("Encountered %d errors during processing:\n", len(errs))
		for _, err := range errs {
			fmt.Printf("  %v\n", err)
		}
	}

	return nil
}

func (a *Analyzer) AnalyzeRecoveriesInConversions(years []int) error {
	var errs []error
	conversionPlaysWithRecoveries := make(map[string][]string, 0)

	for _, year := range years {
		path := common.GetYearGlobPattern(a.inputDir, "nfl", year)
		matches, err := filepath.Glob(path)
		if err != nil {
			fmt.Printf("Error globbing files for year %d: %v\n", year, err)
			continue
		}

		fmt.Printf("year: %d, matches: %v\n", year, len(matches))
		for _, match := range matches {
			result, err := a.processFileNfl(match)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			if len(result.RecoveriesInConversions) > 0 {
				conversionPlaysWithRecoveries[result.ID] = result.RecoveriesInConversions
			}
		}
	}

	// Write results
	if err := a.writeJSONFile("recoveries_in_conversion_plays.json", conversionPlaysWithRecoveries); err != nil {
		return fmt.Errorf("writing recoveries_in_conversion_plays: %w", err)
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
