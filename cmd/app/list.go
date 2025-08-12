package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lpernett/godotenv"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list symbol with data range",
	Long:  "list beginning and ending date ranges for data pertaining to specific ticker and exchange combinations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		configMap, err := godotenv.Read(configFile)
		if err != nil {
			slog.Error("could not read configuration file", "err", err)
			return fmt.Errorf("could not read configurtion file %s", configFile)
		}

		var symbolArray []string
		if len(tickerFile) != 0 {
			symbolArray, err = readTickerFile(tickerFile)
			if err != nil {
				return fmt.Errorf("could not read ticker file %s", err)
			}
		} else {
			symbolArray = args
		}

		conn, err := createEodPsConnection(configMap)
		if err != nil {
			return fmt.Errorf("postgres connection error %s", err)
		}
		defer conn.Close(context.TODO())

		if len(symbolArray) == 0 {
			fmt.Println("list all symbols")
		}

		return nil
	},
}
