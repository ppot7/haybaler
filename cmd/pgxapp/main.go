package main

import (
	"fmt"

	"github.com/ppot7/haybaler/eodps"
)

func main() {
	fmt.Println("Starting PGX Main App")

	priceTable := "DAILY_PRICE_RANGE_VOL"
	dividendTable := "DAILY_DIVIDEND"
	splitTable := "DAILY_SPLIT"

	config := eodps.CreateConnectionConfig("localhost", 5432, "postgres", "Gr8Gaz00!", "DAILY_ASSET_DATA")
	conn, err := eodps.CreateEodPostgresConnection(priceTable, dividendTable, splitTable, config)

	if err != nil {
		fmt.Println(err)
		return
	}

	conn.Close()
	// config := eodps.CreateConnectionConfig(host, port, user, password, database)
	// fmt.Println(config)
	// conn, err := eodps.CreateEodPostgresConnection("price_range_vol", "dividends", "splits", config)
	// if err != nil {
	// 	return
	// }

	// conn.Close()

	fmt.Println("Stopping PGX Main App")
}
