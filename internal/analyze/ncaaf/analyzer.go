package ncaaf

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gamedl/internal/common"
	"gamedl/lib/web/clients/sportsradar"
)

type Analyzer struct {
	inputDir  string
	outputDir string
}

type ProcessResultNcaaf struct {
	ID           string
	Reviews      map[string]int
	BeforeReview map[string][][]string
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

func NewProcessResultNcaaf() ProcessResultNcaaf {
	return ProcessResultNcaaf{
		Reviews:      make(map[string]int),
		BeforeReview: make(map[string][][]string),
	}
}

func (a *Analyzer) processFileNcaaf(path string) (ProcessResultNcaaf, error) {
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

func (a *Analyzer) AnalyzeReviewTypes(years []int) error {
	var errs []error
	typesToGames := make(map[string][]GameReview)
	reviewTypeCount := make(map[string]int)

	for _, year := range years {
		path := common.GetYearGlobPattern(a.inputDir, "ncaaf", year)
		matches, err := filepath.Glob(path)
		if err != nil {
			fmt.Printf("Error globbing files for year %d: %v\n", year, err)
			continue
		}

		fmt.Printf("year: %d, matches: %v\n", year, len(matches))
		for _, match := range matches {
			result, err := a.processFileNcaaf(match)
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

	// Write results
	if err := a.writeJSONFile("types_to_games.json", typesToGames); err != nil {
		return fmt.Errorf("writing types_to_games: %w", err)
	}

	if err := a.writeJSONFile("review_type_count.json", reviewTypeCount); err != nil {
		return fmt.Errorf("writing review_type_count: %w", err)
	}

	// Create review games directory and copy sample games
	reviewGamesDir := filepath.Join(a.outputDir, "review_games")
	if err := os.MkdirAll(reviewGamesDir, 0o755); err != nil {
		fmt.Printf("could not create review_games directory: %v\n", err)
	} else {
		for reviewType, games := range typesToGames {
			if len(games) == 0 {
				continue
			}
			// Take up to 3 most recent games
			nGames := min(len(games), 3)
			lastGames := games[len(games)-nGames:]

			for _, lastGame := range lastGames {
				gameFile := common.GetGameFilePath(a.inputDir, "ncaaf", lastGame.Year, lastGame.ID)
				gameData, err := os.ReadFile(gameFile)
				if err != nil {
					fmt.Printf("could not read game file: %v\n", err)
					continue
				}

				cleanReviewType := reviewType
				if cleanReviewType == "" {
					cleanReviewType = "no_type"
				} else {
					cleanReviewType = strings.ReplaceAll(cleanReviewType, " ", "_")
				}

				baseName := fmt.Sprintf("%s-%s.json", cleanReviewType, lastGame.ID)
				err = os.WriteFile(filepath.Join(reviewGamesDir, baseName), gameData, 0o644)
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
