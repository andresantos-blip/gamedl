package analyze

import (
	"fmt"
	"gamedl/internal/analyze/nba"
	"gamedl/internal/analyze/ncaab"
	"gamedl/internal/analyze/ncaaf"
	"gamedl/internal/analyze/nfl"
	"gamedl/internal/common"
)

type Config struct {
	Competition  string
	AnalysisType string
	InputDir     string
	OutputDir    string
	Seasons      []int
}

func hydrateConfig(config *Config) error {
	// If no years are specified, discover available years from directory structure
	if len(config.Seasons) == 0 {
		availableYears, err := common.GetAvailableYears(config.InputDir, config.Competition)
		if err != nil {
			return fmt.Errorf("failed to discover available years: %w", err)
		}
		if len(availableYears) == 0 {
			return fmt.Errorf("no years found for competition %s in directory %s", config.Competition, config.InputDir)
		}
		config.Seasons = availableYears
		fmt.Printf("Using available years: %v\n", config.Seasons)
	}

	return nil
}

func Run(config Config) error {
	if err := hydrateConfig(&config); err != nil {
		return fmt.Errorf("hydrating config: %w", err)
	}

	switch config.Competition {
	case "nfl":
		return runNFLAnalysis(config)
	case "ncaab":
		return runNCAABAnalysis(config)
	case "ncaaf":
		return runNCAAFAnalysis(config)
	case "nba":
		return runNBAAnalysis(config)
	default:
		return fmt.Errorf("unsupported competition: %s", config.Competition)
	}
}

func runNFLAnalysis(config Config) error {
	analyzer := nfl.NewAnalyzer(config.InputDir, config.OutputDir)

	switch config.AnalysisType {
	case "action-types":
		return analyzer.AnalyzeActionTypes(config.Seasons)
	case "recoveries-in-conversions":
		return analyzer.AnalyzeRecoveriesInConversions(config.Seasons)
	default:
		return fmt.Errorf("unsupported analysis type for NFL: %s", config.AnalysisType)
	}
}

func runNCAABAnalysis(config Config) error {
	analyzer := ncaab.NewAnalyzer(config.InputDir, config.OutputDir)

	switch config.AnalysisType {
	case "review-types":
		return analyzer.AnalyzeReviewTypes(config.Seasons)
	default:
		return fmt.Errorf("unsupported analysis type for NCAAB: %s", config.AnalysisType)
	}
}

func runNCAAFAnalysis(config Config) error {
	analyzer := ncaaf.NewAnalyzer(config.InputDir, config.OutputDir)

	switch config.AnalysisType {
	case "review-types":
		return analyzer.AnalyzeReviewTypes(config.Seasons)
	default:
		return fmt.Errorf("unsupported analysis type for NCAAF: %s", config.AnalysisType)
	}
}

func runNBAAnalysis(config Config) error {
	analyzer := nba.NewAnalyzer(config.InputDir, config.OutputDir)

	switch config.AnalysisType {
	case "lane-violations":
		return analyzer.AnalyzeLaneViolations(config.Seasons)
	default:
		return fmt.Errorf("unsupported analysis type for NBA: %s", config.AnalysisType)
	}
}
