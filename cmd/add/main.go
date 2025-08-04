package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lpernett/godotenv"
	"github.com/ppot7/haybaler"
	"github.com/ppot7/haybaler/eodhdapi"
	"github.com/ppot7/haybaler/eodpostgres"
)

func main() {
	fmt.Println("Starting Add Procedure")

	envMap, err := godotenv.Read("haybaler.config")
	if err != nil {
		fmt.Println(err)
		return
	}

	client := eodhdapi.CreateEodHdClient(envMap["eodhd.host"], envMap["eodhd.token"], nil)

	begin := time.Date(1999, time.November, 1, 0, 0, 0, 0, time.UTC)
	pvData, err := client.RetrieveSplitRecords("MSFT", "US", begin, time.Now())
	if err != nil {
		fmt.Println(err)
		return
	}

	for key, value := range envMap {
		fmt.Println(key, value)
	}

	config, err := eodpostgres.CreateDefaultConfiguration(envMap["eod.ps.host"], envMap["eod.ps.port"], envMap["eod.ps.user"],
		envMap["eod.ps.pwd"], envMap["eod.ps.db"])

	conn, err := eodpostgres.ConnectToPsDatabase(context.TODO(), envMap["eod.ps.schema"], envMap["eod.ps.pv_table"],
		envMap["eod.ps.dividend_table"], envMap["eod.ps.split_table"], config)
	if err != nil {
		fmt.Println("error connecting to postgres db: ", err)
		return
	}
	defer conn.Close(context.TODO())

	dataArray := make([]haybaler.EodSplit, 0, 25)
	for data, err := range pvData {
		if err != nil {
			fmt.Println(err)
		} else {
			dataArray = append(dataArray, *data)
		}
	}

	conn.LoadSplitData(context.TODO(), dataArray)

	fmt.Println("Stopping Add Procedure")
}
