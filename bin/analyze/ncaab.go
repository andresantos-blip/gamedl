package main

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

type ProcessResultNcaab struct {
	ID          string
	EventTypes  map[string]int
	BeforeEvent map[string][][]string
}

func NewProcessResultNcaab() ProcessResultNcaab {
	return ProcessResultNcaab{
		EventTypes:  make(map[string]int),
		BeforeEvent: make(map[string][][]string),
	}
}

var NcaabReviewTypes = []string{
	"challengereview",
	"challengetimeout",
	"requestreview",
	"review",
}

func ProcessFileNcaab(path string) (ProcessResultNcaab, error) {
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

func AnalyzeNcaab() {

	years := []int{2012, 2013, 2014, 2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024}

	var errs []error

	type GameReview struct {
		Year   int        `json:"year"`
		ID     string     `json:"id"`
		Before [][]string `json:"before"`
	}

	eventsToGames := make(map[string][]*GameReview)
	eventTypeCount := make(map[string]int)

	for _, year := range years {
		path := filepath.Join("ncaab_games", strconv.Itoa(year), "*.json")
		matches, _ := filepath.Glob(path)
		fmt.Printf("year: %d, matches: %v \n", year, len(matches))
		for _, match := range matches {
			result, err := ProcessFileNcaab(match)
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

	for t, games := range eventsToGames {
		for i := range games {
			for j := range eventsToGames[t][i].Before {
				start := max(len(eventsToGames[t][i].Before[j])-5, 0)
				eventsToGames[t][i].Before[j] = eventsToGames[t][i].Before[j][start:]
			}
		}
	}

	eventsToGamesData, err := json.MarshalIndent(eventsToGames, "", "  ")
	if err != nil {
		fmt.Printf("could not marshal typesToGames: %v \n", err)
		return
	}

	err = os.WriteFile("review_events_to_games.json", eventsToGamesData, 0644)
	if err != nil {
		fmt.Printf("could not write typesToGames: %v \n", err)
		os.Exit(1)
	}

	eventTypeCountData, err := json.MarshalIndent(eventTypeCount, "", "  ")
	if err != nil {
		fmt.Printf("could not marshal eventTypeCount: %v \n", err)
		os.Exit(1)
	}

	err = os.WriteFile("event_type_count.json", eventTypeCountData, 0644)
	if err != nil {
		fmt.Printf("could not write reviewTypeCount: %v \n", err)
		os.Exit(1)
	}

	for _, err := range errs {
		fmt.Printf("error: %v \n", err)
	}

	err = os.MkdirAll("review_games_ncaab", 0755)
	if err != nil {
		fmt.Printf("could not create review_games directory: %v \n", err)
	}

	for eventType, games := range eventsToGames {
		if len(games) == 0 {
			continue
		}
		// nGames := min(len(games), 3)
		lastGames := games[:]

		for _, lastGame := range lastGames {
			gameFile := filepath.Join("ncaab_games", strconv.Itoa(lastGame.Year), fmt.Sprintf("%s.json", lastGame.ID))
			gameData, err := os.ReadFile(gameFile)
			if err != nil {
				fmt.Printf("could not read game file: %v \n", err)
				continue
			}

			if eventType == "" {
				eventType = "no_type"
			} else {
				eventType = strings.ReplaceAll(eventType, " ", "_")
			}

			baseName := fmt.Sprintf("%s-%s.json", eventType, lastGame.ID)
			err = os.WriteFile(filepath.Join("review_games_ncaab", baseName), gameData, 0644)
			if err != nil {
				fmt.Printf("could not write game file: %v \n", err)
			}
		}

	}
}
