package main

import (
	"fmt"
	"time"

	"github.com/ppot7/haybaler/eodhdapi"
)

func main() {

	begin := time.Date(2003, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := begin.AddDate(0, 3, 0)

	token := "67c86114a60525.79523855"
	client := eodhdapi.CreateEodHdClient("https://www.eodhd.com", token, nil)
	pvRecords, err := client.RetrievePriceVolumeRecords("MSFT", "US", begin, end)
	if err != nil {
		fmt.Println(err)
		return
	}

	for record, err := range pvRecords {
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(record.GoString())
		}
	}

}
