package cmd

import (
	"fmt"
	"gamedl/internal/analyze"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze downloaded game data",
	Long: `Analyze previously downloaded game data files.
Supports various analysis types for different competitions.

Configuration precedence (highest to lowest):
1. Command line flags
2. Environment variables (GAMEDL_*)
3. Configuration file`,
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().StringP("competition", "c", "", "Competition to analyze (nfl, ncaab, ncaaf) (required)")
	analyzeCmd.Flags().StringP("analysis", "a", "", "Analysis type to perform (required)")
	analyzeCmd.Flags().StringP("input-dir", "i", "downloaded_games", "Directory containing downloaded game files")
	analyzeCmd.Flags().StringP("output", "o", "analysis_results", "Output directory for analysis results")
	analyzeCmd.Flags().StringSliceP("years", "y", nil, "Years to include in analysis (e.g., 2023,2024)")

	// Note: We handle required validation in RunE since we use viper for config precedence

	viper.BindPFlag("analyze.competition", analyzeCmd.Flags().Lookup("competition"))
	viper.BindPFlag("analyze.analysis", analyzeCmd.Flags().Lookup("analysis"))
	viper.BindPFlag("analyze.input-dir", analyzeCmd.Flags().Lookup("input-dir"))
	viper.BindPFlag("analyze.output", analyzeCmd.Flags().Lookup("output"))
	viper.BindPFlag("analyze.years", analyzeCmd.Flags().Lookup("years"))

	// Also bind environment variables directly
	viper.BindEnv("analyze.competition", "GAMEDL_ANALYZE_COMPETITION")
	viper.BindEnv("analyze.analysis", "GAMEDL_ANALYZE_ANALYSIS")
	viper.BindEnv("analyze.input-dir", "GAMEDL_ANALYZE_INPUT_DIR")
	viper.BindEnv("analyze.output", "GAMEDL_ANALYZE_OUTPUT")
	viper.BindEnv("analyze.years", "GAMEDL_ANALYZE_YEARS")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	competition := viper.GetString("analyze.competition")
	analysisType := viper.GetString("analyze.analysis")
	inputDir := viper.GetString("analyze.input-dir")
	outputDir := viper.GetString("analyze.output")
	yearsStr := viper.GetStringSlice("analyze.years")

	if competition == "" {
		return fmt.Errorf("competition is required")
	}

	if analysisType == "" {
		return fmt.Errorf("analysis type is required")
	}

	validCompetitions := []string{"nfl", "ncaab", "ncaaf"}

	if !contains(validCompetitions, competition) {
		return fmt.Errorf("invalid competition %s. Valid options: %s", competition, strings.Join(validCompetitions, ", "))
	}

	var years []int
	if len(yearsStr) > 0 {
		for _, y := range yearsStr {
			year, err := parseYear(strings.TrimSpace(y))
			if err != nil {
				return fmt.Errorf("invalid year %s: %w", y, err)
			}
			years = append(years, year)
		}
	}

	fmt.Printf("Analyzing %s data\n", competition)
	fmt.Printf("Analysis type: %s\n", analysisType)
	fmt.Printf("Input directory: %s\n", inputDir)
	fmt.Printf("Output directory: %s\n", outputDir)
	if len(years) > 0 {
		fmt.Printf("Years: %v\n", years)
	} else {
		fmt.Println("Years: all available")
	}

	config := analyze.Config{
		Competition:  competition,
		AnalysisType: analysisType,
		InputDir:     inputDir,
		OutputDir:    outputDir,
		Years:        years,
	}

	if err := analyze.Run(config); err != nil {
		fmt.Fprintf(os.Stderr, "Analysis failed: %v\n", err)
		return err
	}

	return nil
}

func parseYear(s string) (int, error) {
	year := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, fmt.Errorf("invalid year format")
		}
		year = year*10 + int(r-'0')
	}
	if year < 1900 || year > 2100 {
		return 0, fmt.Errorf("year out of valid range")
	}
	return year, nil
}
