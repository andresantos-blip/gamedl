package cmd

import (
	"fmt"
	"os"

	"gamedl/lib/app/build"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "gamedl",
	Short: "A CLI tool for downloading and analyzing game data",
	Long: `GameDL is a command-line interface for downloading game information 
from various sports data providers (SportRadar, BetGenius) and analyzing 
the downloaded data.`,
}

var buildInfo build.Info

func Execute(info build.Info) {
	buildInfo = info
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gamedl.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gamedl")
	}

	viper.SetEnvPrefix("GAMEDL")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
