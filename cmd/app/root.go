package main

import (
	"github.com/spf13/cobra"
)

var (
	tickerFile string

	rootCmd = &cobra.Command{
		Use:   "haybaler",
		Short: "A Data Migration Tool",
		Long: `Haybaler allows the migration of data for
	importing, updating and extracting daily asset data`,
	}

	configFile string
)

func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "f", "haybaler.config", "file containing db and api configurations")
	rootCmd.PersistentFlags().StringVarP(&tickerFile, "ticker-file", "t", "", "file containing ticker.exchange listings")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
}
