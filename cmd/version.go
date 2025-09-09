package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long: `Print the version information for gamedl.
	
Use the --verbose flag to show additional build information including
commit hash and build date.`,
	Run: runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolP("verbose", "v", false, "Show verbose version information including commit and date")
}

func runVersion(cmd *cobra.Command, args []string) {
	verbose, _ := cmd.Flags().GetBool("verbose")

	if verbose {
		fmt.Printf("version: %s\n", buildInfo.Version)
		fmt.Printf("commit: %s\n", buildInfo.Commit)

		// Try to parse and format the date if it's not "unknown"
		if buildInfo.Date != "unknown" {
			if parsedDate, err := time.Parse(time.RFC3339, buildInfo.Date); err == nil {
				fmt.Printf("built: %s\n", parsedDate.Format("2006-01-02 15:04:05 MST"))
			} else {
				fmt.Printf("built: %s\n", buildInfo.Date)
			}
		} else {
			fmt.Printf("built: %s\n", buildInfo.Date)
		}
	} else {
		fmt.Printf("%s\n", buildInfo.Version)
	}
}
