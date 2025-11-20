package cmd

import (
	"fmt"
	"gamedl/internal/download"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download game data from sports providers",
	Long: `Download game information from various sports data providers.
Supports SportRadar and BetGenius providers for different competitions.
`,
	RunE: runDownload,
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringP("competition", "c", "", "Competition to download (values allowed: 'nfl', 'ncaab' or ncaaf) (required)")
	downloadCmd.Flags().StringP("provider", "p", "", "Data provider (values allowed: 'sportradar', 'sr', 'betgenius', 'genius' or 'bg') (required)")
	downloadCmd.Flags().StringSliceP("seasons", "s", nil, "Seasons to download, comma-separated. e.g '2023,2024' (default: all seasons available in the provider)")
	downloadCmd.Flags().IntP("concurrency", "", 10, "Number of concurrent downloads")
	downloadCmd.Flags().StringP("output-dir", "o", "downloaded_games", "Directory to store downloaded game files")

	// Note: We handle required validation in RunE since we use viper for config precedence

	viper.BindPFlag("download.competition", downloadCmd.Flags().Lookup("competition"))
	viper.BindPFlag("download.provider", downloadCmd.Flags().Lookup("provider"))
	viper.BindPFlag("download.seasons", downloadCmd.Flags().Lookup("seasons"))
	viper.BindPFlag("download.concurrency", downloadCmd.Flags().Lookup("concurrency"))
	viper.BindPFlag("download.output-dir", downloadCmd.Flags().Lookup("output-dir"))

	// Also bind environment variables directly
	viper.BindEnv("download.competition", "GAMEDL_DOWNLOAD_COMPETITION")
	viper.BindEnv("download.provider", "GAMEDL_DOWNLOAD_PROVIDER")
	viper.BindEnv("download.seasons", "GAMEDL_DOWNLOAD_SEASONS")
	viper.BindEnv("download.concurrency", "GAMEDL_DOWNLOAD_CONCURRENCY")
	viper.BindEnv("download.output-dir", "GAMEDL_DOWNLOAD_OUTPUT_DIR")
}

func runDownload(cmd *cobra.Command, args []string) error {
	competition := viper.GetString("download.competition")
	provider := viper.GetString("download.provider")
	seasonsStr := viper.GetStringSlice("download.seasons")
	concurrency := viper.GetInt("download.concurrency")
	outputDir := viper.GetString("download.output-dir")

	if competition == "" {
		return fmt.Errorf("competition is required")
	}

	if provider == "" {
		return fmt.Errorf("provider is required")
	}

	validCompetitions := []string{"nfl", "ncaab", "ncaaf", "nba"}
	validProviders := []string{"sportradar", "sr", "betgenius", "genius", "bg"}

	if !contains(validCompetitions, competition) {
		return fmt.Errorf("invalid competition %s. Valid options: %s", competition, strings.Join(validCompetitions, ", "))
	}

	if !contains(validProviders, provider) {
		return fmt.Errorf("invalid provider %s. Valid options: %s", provider, strings.Join(validProviders, ", "))
	}

	var seasons []int
	if len(seasonsStr) > 0 {
		for _, s := range seasonsStr {
			season, err := strconv.Atoi(strings.TrimSpace(s))
			if err != nil {
				return fmt.Errorf("invalid season %s: %w", s, err)
			}
			seasons = append(seasons, season)
		}
	}

	fmt.Printf("Downloading %s data from %s\n", competition, provider)
	if len(seasons) > 0 {
		fmt.Printf("Seasons: %v\n", seasons)
	} else {
		fmt.Println("Seasons: all available")
	}
	fmt.Printf("Concurrency: %d\n", concurrency)
	fmt.Printf("Output directory: %s\n", outputDir)

	config := download.Config{
		Competition: competition,
		Provider:    provider,
		Seasons:     seasons,
		Concurrency: concurrency,
		OutputDir:   outputDir,
	}

	if err := download.Run(config); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		return err
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
