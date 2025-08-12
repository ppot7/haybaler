package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/lpernett/godotenv"
	"github.com/spf13/cobra"
)

var (
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "add ticker.exchange price/dividend/split data to database",
		Long:  "Add historical data for ticker symbol (with exchange) to database",
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

			if len(symbolArray) == 0 {
				slog.Error("no symbols to process")
				return fmt.Errorf("no symbols to process (re- specify arguments or ticker file)")
			}

			client := createEodClient(configMap)
			conn, err := createEodPsConnection(configMap)
			if err != nil {
				return fmt.Errorf("postgres connection error %s", err)
			}
			defer conn.Close(context.TODO())

			begin := time.Date(2003, time.January, 1, 0, 0, 0, 0, time.UTC)
			end := begin.AddDate(1, 0, 0)

			for _, symbol := range symbolArray {
				tickerExchange := strings.SplitN(symbol, ".", 2)

				err := addTicker(context.TODO(), client, conn, tickerExchange[0], tickerExchange[1], begin, end, 50)
				if err != nil {
					slog.Error("add ticker error %s", "err", err)
					return fmt.Errorf("ticker read error: %s", err)
				}

				fmt.Printf("%s --> %s\n", tickerExchange[0], tickerExchange[1])
			}

			return nil
		},
	}
)
