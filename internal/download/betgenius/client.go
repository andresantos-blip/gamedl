package betgenius

import (
	"fmt"
	"gamedl/lib/betgenius"
	"os"
)

func createBetGeniusClient() (*betgenius.Client, error) {
	fixtureKey := os.Getenv("BG_FIXTURE_KEY")
	if fixtureKey == "" {
		return nil, fmt.Errorf("BG_FIXTURE_KEY environment variable not set")
	}

	fixtureUsername := os.Getenv("BG_FIXTURE_USER")
	if fixtureUsername == "" {
		return nil, fmt.Errorf("BG_FIXTURE_USER environment variable not set")
	}

	fixturePassword := os.Getenv("BG_FIXTURE_PASSWORD")
	if fixturePassword == "" {
		return nil, fmt.Errorf("BG_FIXTURE_PASSWORD environment variable not set")
	}

	statsKey := os.Getenv("BG_STATS_KEY")
	if statsKey == "" {
		return nil, fmt.Errorf("BG_STATS_KEY environment variable not set")
	}

	statsUsername := os.Getenv("BG_STATS_USER")
	if statsUsername == "" {
		return nil, fmt.Errorf("BG_STATS_USER environment variable not set")
	}

	statsPassword := os.Getenv("BG_STATS_PASSWORD")
	if statsPassword == "" {
		return nil, fmt.Errorf("BG_STATS_PASSWORD environment variable not set")
	}

	client := betgenius.NewClient(
		betgenius.WithStatsKey(statsKey),
		betgenius.WithFixtureUsername(fixtureUsername),
		betgenius.WithFixturePassword(fixturePassword),
		betgenius.WithFixtureKey(fixtureKey),
		betgenius.WithStatsUsername(statsUsername),
		betgenius.WithStatsPassword(statsPassword),
	)

	return client, nil
}