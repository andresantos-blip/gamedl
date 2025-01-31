package main

import (
	"encoding/json"
	"fmt"
	"gamedl/lib/betgenius"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
)

type ProcessResultNfl struct {
	ID             string
	ActionTypes    map[string]int
	ActionSubTypes map[string]int
	BeforeAction   map[string][][]string
}

func NewProcessResultNfl() ProcessResultNfl {
	return ProcessResultNfl{
		ActionTypes:    make(map[string]int),
		ActionSubTypes: make(map[string]int),
		BeforeAction:   make(map[string][][]string),
	}
}

func ProcessFileNfl(path string) (ProcessResultNfl, error) {
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

func AnalyzeNfl() {

	years := []int{2021, 2022, 2023, 2024}

	var errs []error

	type GameReview struct {
		Year   int        `json:"year"`
		ID     string     `json:"id"`
		Before [][]string `json:"before"`
	}

	actionsToGames := make(map[string][]*GameReview)
	subActionsToGames := make(map[string][]*GameReview)
	actionTypeCount := make(map[string]int)
	subActionTypeCount := make(map[string]int)

	for _, year := range years {
		path := filepath.Join("nfl_games", strconv.Itoa(year), "*.json")
		matches, _ := filepath.Glob(path)
		fmt.Printf("year: %d, matches: %v \n", year, len(matches))
		for _, match := range matches {
			result, err := ProcessFileNfl(match)
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

	actionsToGameData, err := json.MarshalIndent(actionsToGames, "", "  ")
	if err != nil {
		fmt.Printf("could not marshal actionsToGames: %v \n", err)
		return
	}

	err = os.WriteFile("actions_to_games.json", actionsToGameData, 0644)
	if err != nil {
		fmt.Printf("could not write actionsToGameData: %v \n", err)
		os.Exit(1)
	}

	subActionsToGameData, err := json.MarshalIndent(subActionsToGames, "", "  ")
	if err != nil {
		fmt.Printf("could not marshal subActionsToGameData: %v \n", err)
		return
	}

	err = os.WriteFile("sub_actions_to_games.json", subActionsToGameData, 0644)
	if err != nil {
		fmt.Printf("could not write subActionsToGameData: %v \n", err)
		os.Exit(1)
	}

	actionTypeCountData, err := json.MarshalIndent(actionTypeCount, "", "  ")
	if err != nil {
		fmt.Printf("could not marshal actionTypeCount: %v \n", err)
		os.Exit(1)
	}

	err = os.WriteFile("action_type_count.json", actionTypeCountData, 0644)
	if err != nil {
		fmt.Printf("could not write actionTypeCountData: %v \n", err)
		os.Exit(1)
	}

	subActionTypeCountData, err := json.MarshalIndent(subActionTypeCount, "", "  ")
	if err != nil {
		fmt.Printf("could not marshal subActionTypeCount: %v \n", err)
		os.Exit(1)
	}

	err = os.WriteFile("sub_action_type_count.json", subActionTypeCountData, 0644)
	if err != nil {
		fmt.Printf("could not write subActionTypeCountData: %v \n", err)
		os.Exit(1)
	}

	for _, err := range errs {
		fmt.Printf("error: %v \n", err)
	}
	
}
