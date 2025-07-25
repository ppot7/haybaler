package eodhd

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ppot7/haybaler/datautils"
)

type EodHdClient struct {
	Token string
	Url   string
	http.Client
}

func (c *EodHdClient) readRawData(ticker string, exchange string, begin time.Time, end time.Time, flags ...uint32) (io.ReadCloser, error) {
	stmt, err := c.createGetStatement(ticker, exchange, begin, end, flags...)
	if err != nil {
		slog.Error("error creating request. ", "err: ", err)
		return nil, fmt.Errorf("error creating request: %s", err)
	}

	response, err := c.Get(stmt)
	if err != nil {
		slog.Error("error requesting data. ", "err: ", err)
		return nil, fmt.Errorf("error requesting data: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		slog.Error("response error. ", "status code: ", response.Status)
		return nil, fmt.Errorf("response error (Code: %d): %s", response.StatusCode, response.Status)
	}

	slog.Info("response data ready for reading.")

	return response.Body, nil
}

func (c *EodHdClient) RetrievePriceRangeData(ticker string, exchange string, beginDate time.Time, endDate time.Time, flags ...uint32) ([]datautils.PriceRangeData, error) {
	readCloser, err := c.readRawData(ticker, exchange, beginDate, endDate, flags...)
	if err != nil {
		slog.Error("unable to read response")
		return nil, fmt.Errorf("unable to read response: %s", err)
	}
	defer readCloser.Close()

	scanner := bufio.NewScanner(readCloser)
	if !scanner.Scan() {
		slog.Error("no data to return")
		return nil, fmt.Errorf("no data to return")
	}

	dataArray := make([]datautils.PriceRangeData, 0)
	for scanner.Scan() {
		dataPoints := strings.Split(scanner.Text(), ",")
		if scanner.Err() != nil {
			slog.Error("unable to read response line ", "err: ", scanner.Err())
		}

		dataPoint, err := parsePriceRange(ticker, exchange, dataPoints)
		if err != nil {
			slog.Error("unable to parse datapoint")
		}

		dataArray = append(dataArray, *dataPoint)
	}

	return dataArray, nil
}

func (c *EodHdClient) RetrieveDividendData(ticker string, exchange string, beginDate time.Time, endDate time.Time, flags ...uint32) ([]datautils.DividendData, error) {
	flags = append(flags, datautils.DIVIDEND)
	readCloser, err := c.readRawData(ticker, exchange, beginDate, endDate, flags...)
	if err != nil {
		slog.Error("unable to read response")
		return nil, fmt.Errorf("unable to read response: %s", err)
	}
	defer readCloser.Close()

	scanner := bufio.NewScanner(readCloser)
	if !scanner.Scan() {
		slog.Error("no data to return")
		return nil, fmt.Errorf("no data to return")
	}

	dataArray := make([]datautils.DividendData, 0)
	for scanner.Scan() {
		dataPoints := strings.Split(scanner.Text(), ",")
		if scanner.Err() != nil {
			slog.Error("unable to read response line ", "err: ", scanner.Err())
		}

		dataPoint, err := parseDividendData(ticker, exchange, dataPoints)
		if err != nil {
			slog.Error("unable to parse datapoint")
		}
		dataArray = append(dataArray, *dataPoint)
	}

	return dataArray, nil
}

func (c *EodHdClient) RetrieveSplitData(ticker string, exchange string, beginDate time.Time, endDate time.Time, flags ...uint32) ([]datautils.SplitData, error) {
	flags = append(flags, datautils.SPLIT)
	readCloser, err := c.readRawData(ticker, exchange, beginDate, endDate, flags...)
	if err != nil {
		slog.Error("unable to read response")
		return nil, fmt.Errorf("unable to read response: %s", err)
	}
	defer readCloser.Close()

	scanner := bufio.NewScanner(readCloser)
	if !scanner.Scan() {
		slog.Error("no data to return")
		return nil, fmt.Errorf("no data to return")
	}

	dataArray := make([]datautils.SplitData, 0)
	for scanner.Scan() {
		dataPoints := strings.Split(scanner.Text(), ",")
		if scanner.Err() != nil {
			slog.Error("unable to read response line ", "err: ", scanner.Err())
		}

		dataPoint, err := parseSplitData(ticker, exchange, dataPoints)
		if err != nil {
			slog.Error("unable to parse datapoint")
		}
		dataArray = append(dataArray, *dataPoint)
	}

	return dataArray, nil
}

/************ Internal Functions and Types ***********/
func (c *EodHdClient) createGetStatement(ticker string, exchange string, begin time.Time, end time.Time, flags ...uint32) (string, error) {
	stmt := c.Url + "/api/"
	infoType := "eod/"
	formatType := "&fmt=csv"

	if end.Before(begin) {
		slog.Error("ending date before beginning date ", "begin: ", begin.Format(time.DateOnly), " end: ", end.Format(time.DateOnly))
		return "", fmt.Errorf("end date before begin date. begin: %s end: %s", begin.Format(time.DateOnly), end.Format(time.DateOnly))
	}

	for _, flag := range flags {
		switch flag {
		case datautils.DIVIDEND:
			infoType = "div/"
		case datautils.SPLIT:
			infoType = "splits/"
		case datautils.JSON:
			formatType = "&fmt=json"
		}
	}

	stmt += infoType + fmt.Sprintf("%s.%s?api_token=", ticker, exchange) + c.Token
	stmt += fmt.Sprintf("&from=%s", begin.Format(time.DateOnly))
	stmt += fmt.Sprintf("&to=%s", end.Format(time.DateOnly))
	stmt += formatType

	slog.Debug("get request formed. ", "stmt: ", stmt)

	return stmt, nil
}

func parsePriceRange(ticker string, exchange string, records []string) (*datautils.PriceRangeData, error) {

	if len(records) < 7 {
		slog.Error("record array does not contain 5 elements")
		return nil, fmt.Errorf("record array only contains %d elements", len(records))
	}

	tradeDate, err := time.Parse(time.DateOnly, records[0])
	if err != nil {
		slog.Error("error parsing trade date: ", "err", err)
		return nil, fmt.Errorf("error parsing trade date: %v", err)
	}

	open, err := strconv.ParseFloat(records[1], 64)
	if err != nil {
		slog.Error("error parsing open price ", "err: ", err)
		return nil, fmt.Errorf("error parsing open price: %v", err)
	}

	high, err := strconv.ParseFloat(records[2], 64)
	if err != nil {
		slog.Error("error parsing high price ", "err: ", err)
		return nil, fmt.Errorf("error parsing high price: %v", err)
	}

	low, err := strconv.ParseFloat(records[3], 64)
	if err != nil {
		slog.Error("error parsing low price ", "err: ", err)
		return nil, fmt.Errorf("error parsing low price: %v", err)
	}

	close, err := strconv.ParseFloat(records[4], 64)
	if err != nil {
		slog.Error("error parsing close price ", "err: ", err)
		return nil, fmt.Errorf("error parsing close price: %v", err)
	}

	volume, err := strconv.ParseInt(records[6], 10, 64)
	if err != nil {
		slog.Error("error parsing open price ", "err: ", err)
		return nil, fmt.Errorf("error parsing open price: %v", err)
	}

	return &datautils.PriceRangeData{
		Ticker:    ticker,
		Exchange:  exchange,
		TradeDate: tradeDate,
		Open:      open,
		High:      high,
		Low:       low,
		Close:     close,
		Volume:    volume,
	}, nil
}

func parseDividendData(ticker string, exchange string, records []string) (*datautils.DividendData, error) {
	if len(records) < 2 {
		slog.Error("record array does not contain 2 elements")
		return nil, fmt.Errorf("record array only contains %d elements", len(records))
	}

	divDate, err := time.Parse(time.DateOnly, records[0])
	if err != nil {
		slog.Error("error parsing ex-dividend date: ", "err", err)
		return nil, fmt.Errorf("error parsing ex-dividend date: %v", err)
	}

	divValue, err := strconv.ParseFloat(records[1], 64)
	if err != nil {
		slog.Error("error parsing dividend value ", "err: ", err)
		return nil, fmt.Errorf("error parsing dividend value: %v", err)
	}

	return &datautils.DividendData{
		Ticker:        ticker,
		Exchange:      exchange,
		ExDate:        divDate,
		DividendValue: divValue,
	}, nil
}

func parseSplitData(ticker string, exchange string, records []string) (*datautils.SplitData, error) {
	if len(records) < 2 {
		slog.Error("record array does not contain 2 elements")
		return nil, fmt.Errorf("record array only contains %d elements", len(records))
	}

	splitDate, err := time.Parse(time.DateOnly, records[0])
	if err != nil {
		slog.Error("error parsing split date: ", "err", err)
		return nil, fmt.Errorf("error parsing split date: %v", err)
	}

	splitFactor := strings.Split(records[1], "/")
	dividend, err := strconv.ParseFloat(splitFactor[0], 64)
	if err != nil {
		slog.Error("could not compute split factor ", "err", err)
		return nil, fmt.Errorf("could not compute split factor: %s", err)
	}

	divisor, err := strconv.ParseFloat(splitFactor[1], 64)
	if err != nil || divisor == 0.0 {
		slog.Error("could not compute split factor ", "err", err)
		return nil, fmt.Errorf("could not compute split factor: %s", err)
	}

	return &datautils.SplitData{
		Ticker:      ticker,
		Exchange:    exchange,
		SplitDate:   splitDate,
		SplitFactor: dividend / divisor,
	}, nil

}
