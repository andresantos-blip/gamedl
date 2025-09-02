package ncaab

import "fmt"

type Analyzer struct {
	outputDir string
}

func NewAnalyzer(outputDir string) *Analyzer {
	return &Analyzer{
		outputDir: outputDir,
	}
}

func (a *Analyzer) AnalyzeActions(years []int) error {
	fmt.Println("NCAAB analysis not yet implemented")
	return fmt.Errorf("NCAAB analysis not yet implemented")
}