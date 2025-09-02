package analyze

import (
	"fmt"
	"gamedl/internal/analyze/ncaab"
	"gamedl/internal/analyze/ncaaf"
	"gamedl/internal/analyze/nfl"
)

type Config struct {
	Competition  string
	AnalysisType string
	InputDir     string
	OutputDir    string
	Years        []int
}

func Run(config Config) error {
	switch config.Competition {
	case "nfl":
		return runNFLAnalysis(config)
	case "ncaab":
		return runNCAABAnalysis(config)
	case "ncaaf":
		return runNCAAFAnalysis(config)
	default:
		return fmt.Errorf("unsupported competition: %s", config.Competition)
	}
}

func runNFLAnalysis(config Config) error {
	analyzer := nfl.NewAnalyzer(config.InputDir, config.OutputDir)

	switch config.AnalysisType {
	case "action-types":
		return analyzer.AnalyzeActionTypes(config.Years)
	default:
		return fmt.Errorf("unsupported analysis type for NFL: %s", config.AnalysisType)
	}
}

func runNCAABAnalysis(config Config) error {
	analyzer := ncaab.NewAnalyzer(config.InputDir, config.OutputDir)

	switch config.AnalysisType {
	case "review-types":
		return analyzer.AnalyzeReviewTypes(config.Years)
	default:
		return fmt.Errorf("unsupported analysis type for NCAAB: %s", config.AnalysisType)
	}
}

func runNCAAFAnalysis(config Config) error {
	analyzer := ncaaf.NewAnalyzer(config.InputDir, config.OutputDir)

	switch config.AnalysisType {
	case "review-types":
		return analyzer.AnalyzeReviewTypes(config.Years)
	default:
		return fmt.Errorf("unsupported analysis type for NCAAF: %s", config.AnalysisType)
	}
}
