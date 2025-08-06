package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lpernett/godotenv"
	"github.com/ppot7/haybaler"
	"github.com/spf13/cobra"
)

var (
	tickerFile string
	addCmd     = &cobra.Command{
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
					slog.Error("add ticker error %s", err)
				}

				fmt.Printf("%s --> %s\n", tickerExchange[0], tickerExchange[1])
			}

			return nil
		},
	}
)

func init() {
	addCmd.PersistentFlags().StringVarP(&tickerFile, "ticker-file", "t", "", "file containing ticker.exchange listings")
}

func readTickerFile(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		slog.Error("could not open ticker file", "err", err)
		return nil, fmt.Errorf("could not open ticker file: %s", fileName)
	}
	defer file.Close()

	var symbol string
	tickerArray := make([]string, 0, 20)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		symbol = scanner.Text()
		if !strings.HasPrefix(symbol, "#") {
			tickerArray = append(tickerArray, symbol)
		}
	}

	return tickerArray, nil
}

func addTicker(ctx context.Context, client haybaler.EodDataRetriever, conn haybaler.EodDataLoader, ticker string, exchange string, begin time.Time, end time.Time, batchSize int) error {

	pvRecords, err := client.RetrievePriceVolumeData(ticker, exchange, begin, end)
	if err != nil {
		slog.Error("price/volumevretrieval error", "err", err)
	}

	err = conn.LoadPriceVolumeStream(ctx, pvRecords, batchSize)
	if err != nil {
		slog.Error("error loading price/volume data", "err", err)
	}

	divRecords, err := client.RetrieveDividendData(ticker, exchange, begin, end)
	if err != nil {
		slog.Error("dividend retrieval error", "err", err)
	}

	err = conn.LoadDividendStream(ctx, divRecords, batchSize)
	if err != nil {
		slog.Error("error loading dividend data", "err", err)
	}

	splitRecords, err := client.RetrieveSplitData(ticker, exchange, begin, end)
	if err != nil {
		slog.Error("split retrieval error", "err", err)
	}

	err = conn.LoadSplitStream(ctx, splitRecords, batchSize)
	if err != nil {
		slog.Error("error loading split data", "err", err)
	}

	return nil
}
