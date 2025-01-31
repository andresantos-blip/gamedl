package main

import (
	"encoding/json"
	"fmt"
	"gamedl/lib/sportsradar"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ProcessResultNcaaf struct {
	ID           string
	Reviews      map[string]int
	BeforeReview map[string][][]string
}

func NewProcessResultNcaaf() ProcessResultNcaaf {
	return ProcessResultNcaaf{
		Reviews:      make(map[string]int),
		BeforeReview: make(map[string][][]string),
	}
}

func ProcessFileNcaaf(path string) (ProcessResultNcaaf, error) {
	result := NewProcessResultNcaaf()
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return result, fmt.Errorf("could not open file %s: %w", path, err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return result, fmt.Errorf("could not read file %s: %w", path, err)
	}

	pbpData := &sportsradar.NcaafGamePbp{}
	err = json.Unmarshal(data, pbpData)
	if err != nil {
		return result, fmt.Errorf("could not unmarshal game pbp: %w", err)
	}

	periods := pbpData.Periods
	for _, period := range periods {
		pbpEvents := period.Pbp
		for _, pbpEvent := range pbpEvents {
			for _, detailedEvent := range pbpEvent.Events {
				for i, eventDetails := range detailedEvent.Details {
					if eventDetails.Review != nil {
						review := eventDetails.Review
						if review.Result == "overturned" {
							result.Reviews[review.Type]++
							beforeReview := make([]string, 0, i)
							for j := 0; j < i; j++ {
								beforeReview = append(beforeReview, detailedEvent.Details[j].Category)
							}
							result.BeforeReview[review.Type] = append(result.BeforeReview[review.Type], beforeReview)
						}
					}
				}
			}
		}
	}
	result.ID = pbpData.ID
	return result, nil
}

func AnalyzeNcaaf() {

	years := []int{2013, 2014, 2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024}

	var errs []error

	type GameReview struct {
		Year   int        `json:"year"`
		ID     string     `json:"id"`
		Before [][]string `json:"before"`
	}

	typesToGames := make(map[string][]GameReview)
	reviewTypeCount := make(map[string]int)

	for _, year := range years {
		path := filepath.Join("ncaaf_games", strconv.Itoa(year), "*.json")
		matches, _ := filepath.Glob(path)
		fmt.Printf("year: %d, matches: %v \n", year, len(matches))
		for _, match := range matches {
			result, err := ProcessFileNcaaf(match)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			for reviewType, count := range result.Reviews {
				typesToGames[reviewType] = append(
					typesToGames[reviewType],
					GameReview{Year: year, ID: result.ID, Before: result.BeforeReview[reviewType]},
				)
				reviewTypeCount[reviewType] += count
			}

		}
	}

	typesToGamesData, err := json.Marshal(typesToGames)
	if err != nil {
		fmt.Printf("could not marshal typesToGames: %v \n", err)
		return
	}

	err = os.WriteFile("types_to_games.json", typesToGamesData, 0644)
	if err != nil {
		fmt.Printf("could not write typesToGames: %v \n", err)
		os.Exit(1)
	}

	reviewTypeCountData, err := json.Marshal(reviewTypeCount)
	if err != nil {
		fmt.Printf("could not marshal reviewTypeCount: %v \n", err)
		os.Exit(1)
	}

	err = os.WriteFile("review_type_count.json", reviewTypeCountData, 0644)
	if err != nil {
		fmt.Printf("could not write reviewTypeCount: %v \n", err)
		os.Exit(1)
	}

	for _, err := range errs {
		fmt.Printf("error: %v \n", err)
	}

	for t, games := range typesToGames {
		if len(games) == 0 {
			continue
		}
		nGames := min(len(games), 3)
		lastGames := games[len(games)-nGames:]

		for _, lastGame := range lastGames {
			gameFile := filepath.Join("ncaaf_games", strconv.Itoa(lastGame.Year), fmt.Sprintf("%s.json", lastGame.ID))
			gameData, err := os.ReadFile(gameFile)
			if err != nil {
				fmt.Printf("could not read game file: %v \n", err)
				continue
			}

			if t == "" {
				t = "no_type"
			} else {
				t = strings.ReplaceAll(t, " ", "_")
			}

			err = os.MkdirAll("review_games", 0755)
			if err != nil {
				fmt.Printf("could not create review_games directory: %v \n", err)
				continue
			}
			baseName := fmt.Sprintf("%s-%s.json", t, lastGame.ID)
			err = os.WriteFile(filepath.Join("review_games", baseName), gameData, 0644)
			if err != nil {
				fmt.Printf("could not write game file: %v \n", err)
			}
		}

	}
}
