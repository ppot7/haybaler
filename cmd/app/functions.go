package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/ppot7/haybaler"
	"github.com/ppot7/haybaler/eodhdapi"
	"github.com/ppot7/haybaler/eodpostgres"
)

func createEodClient(configMap map[string]string) *eodhdapi.EodHdApiClient {
	return eodhdapi.CreateEodHdClient(configMap["eodhd.host"], configMap["eodhd.token"], nil)
}

func createEodPsConnection(configMap map[string]string) (*eodpostgres.EodPsConnection, error) {
	config, err := eodpostgres.CreateDefaultConfiguration(configMap["eod.ps.host"], configMap["eod.ps.port"],
		configMap["eod.ps.user"], configMap["eod.ps.pwd"], configMap["eod.ps.db"])
	if err != nil {
		slog.Error("could not establish connection to ps database", "err", err)
		return nil, fmt.Errorf("could not establish connection to ps database")
	}

	conn, err := eodpostgres.ConnectToPsDatabase(context.TODO(), configMap["eod.ps.schema"], configMap["eod.ps.pv_table"],
		configMap["eod.ps.dividend_table"], configMap["eod.ps.split_table"], config)
	if err != nil {
		slog.Error("connection error ", "err", err)
		return nil, fmt.Errorf("connection error %s", err)
	}

	return conn, nil
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
