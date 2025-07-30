package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ppot7/haybaler/eodhd"
	"github.com/ppot7/haybaler/eodps"
)

func main() {
	fmt.Println("Starting PGX Main App")

	token := "67c86114a60525.79523855"

	eodClient := &eodhd.EodHdClient{
		Token:  token,
		Url:    "https://www.eodhd.com",
		Client: http.Client{},
	}

	begin := time.Date(2003, time.January, 1, 0, 0, 0, 0, time.UTC)
	dataPoints, err := eodClient.RetrievePriceRangeData("MSFT", "US", begin, begin.AddDate(0, 1, 0))
	if err != nil {
		slog.Error("error retrieving data points")
		return
	}

	pgConn, err := eodps.CreateDefaultEodConnection("localhost", 5432, "postgres", "Gr8Gaz00!", "DAILY_ASSET_DATA")
	if err != nil {
		slog.Error("pgx connection error", "err", err)
		return
	}
	defer pgConn.Close()

	pgConn.IngestPriceRangeData(context.TODO(), dataPoints)

	for _, dataPoint := range dataPoints {
		fmt.Printf("%s\n", dataPoint.GoString())
	}

	fmt.Println("Stopping PGX Main App")
}
