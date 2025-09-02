package analyze

import (
	"fmt"
	"gamedl/internal/analyze/nfl"
	"gamedl/internal/analyze/ncaab"
	"gamedl/internal/analyze/ncaaf"
)

type Config struct {
	Competition  string
	AnalysisType string
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
	analyzer := nfl.NewAnalyzer(config.OutputDir)
	
	switch config.AnalysisType {
	case "actions":
		return analyzer.AnalyzeActions(config.Years)
	default:
		return fmt.Errorf("unsupported analysis type for NFL: %s", config.AnalysisType)
	}
}

func runNCAABAnalysis(config Config) error {
	analyzer := ncaab.NewAnalyzer(config.OutputDir)
	
	switch config.AnalysisType {
	case "actions":
		return analyzer.AnalyzeActions(config.Years)
	default:
		return fmt.Errorf("unsupported analysis type for NCAAB: %s", config.AnalysisType)
	}
}

func runNCAAFAnalysis(config Config) error {
	analyzer := ncaaf.NewAnalyzer(config.OutputDir)
	
	switch config.AnalysisType {
	case "actions":
		return analyzer.AnalyzeActions(config.Years)
	default:
		return fmt.Errorf("unsupported analysis type for NCAAF: %s", config.AnalysisType)
	}
}