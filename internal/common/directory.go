package common

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

// GetGamesDirectoryName returns the directory name for a given competition
func GetGamesDirectoryName(competition string) string {
	return competition
}

// GetGamesDirectoryPath returns the full path to the games directory for a competition
func GetGamesDirectoryPath(baseDir, competition string) string {
	return filepath.Join(baseDir, GetGamesDirectoryName(competition))
}

// GetYearDirectoryPath returns the full path to the year directory within a competition's games directory
func GetYearDirectoryPath(baseDir, competition string, year int) string {
	return filepath.Join(GetGamesDirectoryPath(baseDir, competition), strconv.Itoa(year))
}

// GetGameFilePath returns the full path to a specific game file
func GetGameFilePath(baseDir, competition string, year int, gameID string) string {
	return filepath.Join(GetYearDirectoryPath(baseDir, competition, year), gameID+".json")
}

// GetYearGlobPattern returns the glob pattern for all game files in a specific year
func GetYearGlobPattern(baseDir, competition string, year int) string {
	return filepath.Join(GetYearDirectoryPath(baseDir, competition, year), "*.json")
}

// CreateYearDirectory creates the year directory for a competition if it doesn't exist
func CreateYearDirectory(baseDir, competition string, year int) error {
	yearDir := GetYearDirectoryPath(baseDir, competition, year)
	return os.MkdirAll(yearDir, 0o755)
}

// GetAvailableYears returns all available years for a competition by examining the directory structure
func GetAvailableYears(baseDir, competition string) ([]int, error) {
	gamesDir := GetGamesDirectoryPath(baseDir, competition)

	// Check if the games directory exists
	if _, err := os.Stat(gamesDir); os.IsNotExist(err) {
		return []int{}, fmt.Errorf("games directory does not exist: %s", gamesDir)
	}

	// Read the directory contents
	entries, err := os.ReadDir(gamesDir)
	if err != nil {
		return []int{}, fmt.Errorf("failed to read games directory %s: %w", gamesDir, err)
	}
	currentYear := time.Now().Year()

	var years []int
	for _, entry := range entries {
		// Only process directories
		if !entry.IsDir() {
			continue
		}

		// Try to parse the directory name as a year
		yearStr := entry.Name()
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			// Skip directories that don't represent years
			continue
		}

		// Validate year range (reasonable bounds)
		if year >= 1900 && year <= (currentYear+1) {
			years = append(years, year)
		}
	}

	// Sort years in ascending order
	sort.Ints(years)

	return years, nil
}
