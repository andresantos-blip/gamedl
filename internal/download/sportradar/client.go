package sportradar

import (
	"fmt"
	"gamedl/lib/sportsradar"
	"os"
)

func createSportRadarClientWithNCAB() (*sportsradar.Client, error) {
	apiKey := os.Getenv("SPORTRADAR_NCAAB_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("SPORTRADAR_NCAAB_KEY environment variable not set")
	}

	client := sportsradar.NewClient(sportsradar.WithNcaabAPIKey(apiKey))
	return client, nil
}

func createSportRadarClientWithNCAF() (*sportsradar.Client, error) {
	apiKey := os.Getenv("SPORTRADAR_NCAAF_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("SPORTRADAR_NCAAF_KEY environment variable not set")
	}

	client := sportsradar.NewClient(sportsradar.WithNcaafAPIKey(apiKey))
	return client, nil
}