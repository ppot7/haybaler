package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lpernett/godotenv"
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

	config, err := eodpostgres.CreateDefaultConfiguration(envMap["eod.ps.host"], envMap["eod.ps.port"], envMap["eod.ps.user"],
		envMap["eod.ps.pwd"], envMap["eod.ps.db"])
	if err != nil {
		fmt.Println("error connecting to postgres db: ", err)
		return
	}
	conn, err := eodpostgres.ConnectToPsDatabase(context.TODO(), envMap["eod.ps.schema"], envMap["eod.ps.pv_table"],
		envMap["eod.ps.dividend_table"], envMap["eod.ps.split_table"], config)
	if err != nil {
		fmt.Println("error connecting to postgres db: ", err)
		return
	}
	defer conn.Close(context.TODO())

	begin := time.Date(1999, time.November, 1, 0, 0, 0, 0, time.UTC)

	// pvData, err := client.RetrieveSplitData("MSFT", "US", begin, time.Now())
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// conn.LoadSplitStream(context.TODO(), pvData, 25)

	divData, err := client.RetrieveDividendData("MSFT", "US", begin, time.Now())
	if err != nil {
		fmt.Println(err)
	}

	conn.LoadDividendStream(context.TODO(), divData, 25)

	fmt.Println("Stopping Add Procedure")
}
