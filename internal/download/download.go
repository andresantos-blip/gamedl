package download

import (
	"fmt"
	"gamedl/internal/download/betgenius"
	"gamedl/internal/download/sportradar"
)

type Config struct {
	Competition string
	Provider    string
	Seasons     []int
	Concurrency int
	OutputDir   string
}

func Run(config Config) error {
	switch config.Provider {
	case "betgenius":
		return runBetGenius(config)
	case "sportradar":
		return runSportRadar(config)
	default:
		return fmt.Errorf("unsupported provider: %s", config.Provider)
	}
}

func runBetGenius(config Config) error {
	switch config.Competition {
	case "nfl":
		return betgenius.DownloadNFL(config.Seasons, config.Concurrency, config.OutputDir)
	case "ncaab":
		return betgenius.DownloadNCAB(config.Seasons, config.Concurrency, config.OutputDir)
	case "ncaaf":
		return betgenius.DownloadNCAF(config.Seasons, config.Concurrency, config.OutputDir)
	default:
		return fmt.Errorf("unsupported competition for BetGenius: %s", config.Competition)
	}
}

func runSportRadar(config Config) error {
	switch config.Competition {
	case "nfl":
		return sportradar.DownloadNFL(config.Seasons, config.Concurrency, config.OutputDir)
	case "ncaab":
		return sportradar.DownloadNCAB(config.Seasons, config.Concurrency, config.OutputDir)
	case "ncaaf":
		return sportradar.DownloadNCAF(config.Seasons, config.Concurrency, config.OutputDir)
	default:
		return fmt.Errorf("unsupported competition for SportRadar: %s", config.Competition)
	}
}
