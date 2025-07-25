package eodhd

import (
	"net/http"
	"testing"
	"time"

	"github.com/ppot7/haybaler/datautils"
)

func TestCreateGetStatement(t *testing.T) {

	client := &EodHdClient{
		Token:  "MyToken",
		Url:    "https://www.eodhd.com",
		Client: http.Client{},
	}

	begin := time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)
	end := begin.AddDate(0, 0, 7)
	stmt, err := client.createGetStatement("MSFT", "US", begin, end, datautils.PRICE)
	answer := "https://www.eodhd.com/api/eod/MSFT.US?api_token=MyToken&from=2023-12-01&to=2023-12-08"
	if err != nil {
		t.Error(err)
	}
	if stmt != answer {
		t.Errorf("Get statements do not match: %s %s", stmt, answer)
	}

	stmt, err = client.createGetStatement("MSFT", "US", begin, end, datautils.DIVIDEND)
	if err != nil || stmt == answer {
		t.Errorf("statement should not match: %s", stmt)
	}
	end = end.AddDate(0, 0, -14)
	_, err = client.createGetStatement("MSFT", "US", begin, end)
	if err == nil {
		t.Errorf("error not triggered begin: %s end: %s", begin.Format(time.DateOnly), end.Format(time.DateOnly))
	}

}
