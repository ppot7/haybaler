package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lpernett/godotenv"
	"github.com/ppot7/haybaler"
	"github.com/ppot7/haybaler/eodhdapi"
)

func main() {
	fmt.Println("Starting Add Procedure")

	envMap, err := godotenv.Read("haybaler.config")
	if err != nil {
		fmt.Println(err)
		return
	}

	client := eodhdapi.CreateEodHdClient(envMap["eodhd.host"], envMap["eodhd.token"], nil)

	// config, err := eodpostgres.CreateDefaultConfiguration(envMap["eod.ps.host"], envMap["eod.ps.port"], envMap["eod.ps.user"],
	// 	envMap["eod.ps.pwd"], envMap["eod.ps.db"])
	// if err != nil {
	// 	fmt.Println("error connecting to postgres db: ", err)
	// 	return
	// }
	// conn, err := eodpostgres.ConnectToPsDatabase(context.TODO(), envMap["eod.ps.schema"], envMap["eod.ps.pv_table"],
	// 	envMap["eod.ps.dividend_table"], envMap["eod.ps.split_table"], config)
	// if err != nil {
	// 	fmt.Println("error connecting to postgres db: ", err)
	// 	return
	// }
	// defer conn.Close(context.TODO())

	begin := time.Date(1999, time.November, 1, 0, 0, 0, 0, time.UTC)

	AddTicker(context.TODO(), client, "MSFT", "US", begin, time.Now())

	pvData, err := client.RetrieveSplitData("MSFT", "US", begin, time.Now())
	if err != nil {
		fmt.Println(err)
		return
	}

	dataArray := make([]haybaler.EodSplit, 0, 25)
	for data, err := range pvData {
		if err != nil {
			fmt.Println(err)
		} else {
			dataArray = append(dataArray, *data)
		}
	}

	// conn.LoadSplitData(context.TODO(), dataArray)

	fmt.Println("Stopping Add Procedure")
}

func AddTicker(ctx context.Context, retriever haybaler.EodDataRetriever, ticker string, exchange string, begin time.Time, end time.Time) error {

	pvData, err := retriever.RetrieveSplitData(ticker, exchange, begin, end)
	if err != nil {
		return fmt.Errorf("error retrieving data from api: %s", err)
	}

	batchSize := 25
	batchCount := 0
	batchData := make([]haybaler.EodSplit, 0, batchSize)
	for data := range pvData {
		batchCount++
		batchData = append(batchData, *data)
		if batchCount%batchSize == 0 {
			fmt.Printf("Batch Size: %d", len(batchData))
			batchData = make([]haybaler.EodSplit, 0, batchSize)
			batchCount = 0
		}
	}

	if len(batchData) > 0 {
		fmt.Printf("Batch Size: %d", len(batchData))
	}

	return nil
}
