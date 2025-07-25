package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ppot7/haybaler/eodhd"
)

func main() {
	fmt.Println("Starting Haybaler")
	token := "67c86114a60525.79523855"

	eodClient := &eodhd.EodHdClient{
		Token:  token,
		Url:    "https://www.eodhd.com",
		Client: http.Client{},
	}

	begin := time.Date(2003, time.February, 1, 0, 0, 0, 0, time.UTC)
	dataPoints, err := eodClient.RetrieveSplitData("MSFT", "US", begin, begin.AddDate(20, 3, 0))
	if err != nil {
		slog.Error("error retrieving data points")
		return
	}

	for _, dataPoint := range dataPoints {
		fmt.Printf("%s\n", dataPoint.GoString())
	}

	// "QUESTDB_URL": "http://localhost:9000",
	//             "QUESTDB_USER": "questdb",
	//             "QUESTDB_PASSWORD": "questdb"
	//         }

	fmt.Println("Stopping Haybaler")
}
