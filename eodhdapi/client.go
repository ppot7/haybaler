package eodhdapi

import (
	"bufio"
	"fmt"
	"iter"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ppot7/haybaler"
)

const (
	dividend = 0b00000000000000010000000000000000
	split    = 0b00000000000000100000000000000000
	price    = 0b00000000000001000000000000000000
)

type EodHdApiClient struct {
	haybaler.EodHttpClient
	token string
}

func CreateEodHdClient(host string, token string, client *http.Client) *EodHdApiClient {
	/* if client is not specified use http.Client as default */
	if client == nil {
		client = &http.Client{}
	}

	return &EodHdApiClient{
		EodHttpClient: haybaler.EodHttpClient{
			Url:    host,
			Client: *client,
		},
		token: token,
	}
}

func (p *EodHdApiClient) RetrievePriceVolumeRecords(ticker string, exchange string, begin time.Time, end time.Time, options ...uint32) (iter.Seq2[*haybaler.EodPriceVolume, error], error) {
	rawSet, err := p.requestRawRecords(ticker, exchange, begin, end, options...)
	if err != nil {
		slog.Error("could not retrieve raw data", "err", err)
		return nil, fmt.Errorf("could not retrieve raw data (error: %s)", err)
	}

	return func(yield func(*haybaler.EodPriceVolume, error) bool) {
		for rawRecord := range rawSet {
			record, err := parsePriceVolumeCsv(ticker, exchange, rawRecord)
			if !yield(record, err) {
				return
			}
		}
	}, nil
}

/************** Internal Functions ******************/
func (p *EodHdApiClient) requestRawRecords(ticker string, exchange string, begin time.Time, end time.Time, options ...uint32) (iter.Seq[string], error) {

	getStmt, err := p.buildGetStatement(ticker, exchange, begin, end, options...)
	if err != nil {
		return nil, fmt.Errorf("could not form proper Get request: %s", err)
	}

	response, err := p.Get(getStmt)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve response: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		slog.Error("status code error", "status", response.Status)
		return nil, fmt.Errorf("invalid response status: %s", response.Status)
	}

	return func(yield func(string) bool) {
		scanner := bufio.NewScanner(response.Body)
		defer response.Body.Close()

		scanner.Scan()
		for scanner.Scan() {
			if !yield(scanner.Text()) {
				return
			}
		}
	}, nil
}

func (p *EodHdApiClient) buildGetStatement(ticker string, exchange string, begin time.Time, end time.Time, options ...uint32) (string, error) {
	infoType := "eod"

	if end.Before(begin) {
		slog.Error("end date precedes begin date", "begin", begin.Format(time.DateOnly), "end", end.Format(time.DateOnly))
		return "", fmt.Errorf("end date (%s) precedes begin date (%s)", end.Format(time.DateOnly), begin.Format(time.DateOnly))
	}

	for _, flag := range options {
		switch flag {
		case dividend:
			infoType = "div"
		case split:
			infoType = "splits"
		}
	}

	stmt := fmt.Sprintf("%s/api/%s/%s.%s?api_token=%s&from=%s&to=%s", p.Url, infoType, ticker, exchange, p.token, begin.Format(time.DateOnly), end.Format(time.DateOnly))
	slog.Debug("Url formation successful", "stmt", stmt)

	return stmt, nil
}

func parsePriceVolumeCsv(ticker string, exchange string, dataCsv string) (*haybaler.EodPriceVolume, error) {
	/* Parse .csv line into 6 elements - TradeDate, Ticker, Exchange, Open, High, Low, Close, Volume */
	dataPoints := strings.Split(dataCsv, ",")
	if len(dataPoints) != 7 {
		slog.Error("data .csv string does not contain 6 elements", "data", dataCsv, "len", len(dataPoints))
		return nil, fmt.Errorf("cannot parse data string: %s", dataCsv)
	}
	/* Parse Element #1 - TradeDate */
	tradeDate, err := time.Parse("2006-01-02", dataPoints[0])
	if err != nil {
		slog.Error("date parse error", "date str", dataPoints[0])
		return nil, fmt.Errorf("error parsing trade date string (%s)", dataPoints[0])
	}
	/* Parse Element #2 - Open Price */
	open, err := strconv.ParseFloat(dataPoints[1], 32)
	if err != nil {
		slog.Error("open price parse error", "input", dataPoints[1])
		return nil, fmt.Errorf("error parsing open price string (%s) to float32", dataPoints[1])
	}
	/* Parse Element #3 - High Price */
	high, err := strconv.ParseFloat(dataPoints[2], 32)
	if err != nil {
		slog.Error("high price parse error", "input", dataPoints[2])
		return nil, fmt.Errorf("error parsing high price string (%s) to float32", dataPoints[2])
	}
	/* Parse Element #4 - Low Price */
	low, err := strconv.ParseFloat(dataPoints[3], 32)
	if err != nil {
		slog.Error("low price parse error", "input", dataPoints[3])
		return nil, fmt.Errorf("error parsing low price string (%s) to float32", dataPoints[3])
	}
	/* Parse Element #5 - Close Price */
	close, err := strconv.ParseFloat(dataPoints[4], 32)
	if err != nil {
		slog.Error("close price parse error", "input", dataPoints[4])
		return nil, fmt.Errorf("error parsing close price string (%s) to float32", dataPoints[4])
	}
	/* Parse Element #6 -  Volume */
	volume, err := strconv.ParseInt(dataPoints[6], 10, 32)
	if err != nil {
		slog.Error("volume parse error", "input", dataPoints[6])
		return nil, fmt.Errorf("error parsing volume string (%s) to int32", dataPoints[6])
	}
	return &haybaler.EodPriceVolume{
		TradeDate: tradeDate,
		Ticker:    ticker,
		Exchange:  exchange,
		Open:      float32(open),
		High:      float32(high),
		Low:       float32(low),
		Close:     float32(close),
		Volume:    int32(volume),
	}, nil
}

func parseDividendCsv(ticker string, exchange string, dataCsv string) (*haybaler.EodDividend, error) {
	/* Parse .csv line into 6 elements - ExDate, Ticker, Exchange, Dividend */
	dataPoints := strings.Split(dataCsv, ",")
	if len(dataPoints) != 2 {
		slog.Error("data .csv string does not contain 2 elements", "data", dataCsv)
		return nil, fmt.Errorf("cannot parse data string: %s", dataCsv)
	}
	/* Parse Element #1 - ExDate */
	exDate, err := time.Parse("2006-01-02", dataPoints[0])
	if err != nil {
		slog.Error("date parse error", "date str", dataPoints[0])
		return nil, fmt.Errorf("error parsing ex-dividend date string (%s)", dataPoints[0])
	}
	/* Parse Element #2 - Dividend Value */
	dividend, err := strconv.ParseFloat(dataPoints[1], 64)
	if err != nil {
		slog.Error("error parsing dividend value ", "string ", dataPoints[1])
		return nil, fmt.Errorf("error parsing dividend value string (%s)", dataPoints[1])
	}

	return &haybaler.EodDividend{
		ExDate:   exDate,
		Ticker:   ticker,
		Exchange: exchange,
		Value:    float32(dividend),
	}, nil
}

func parseSplitCsv(ticker string, exchange string, dataCsv string) (*haybaler.EodSplit, error) {
	/* Parse .csv line into 6 elements - ExDate, Ticker, Exchange, Dividend */
	dataPoints := strings.Split(dataCsv, ",")
	if len(dataPoints) != 2 {
		slog.Error("data .csv string does not contain 2 elements", "data", dataCsv)
		return nil, fmt.Errorf("cannot parse data string: %s", dataCsv)
	}
	/* Parse Element #1 - ExDate */
	splitDate, err := time.Parse("2006-01-02", dataPoints[0])
	if err != nil {
		slog.Error("date parse error", "date str", dataPoints[0])
		return nil, fmt.Errorf("error parsing split date string (%s)", dataPoints[0])
	}
	/* Parse Element #2 - Split Factor */
	splitValues := strings.Split(dataPoints[1], "/")
	if len(splitValues) != 2 {
		slog.Error("split factor string does not contain 2 elements", "data", dataPoints[1])
		return nil, fmt.Errorf("cannot parse split factor string: %s", dataPoints[1])
	}

	dividend, err := strconv.ParseFloat(splitValues[0], 64)
	if err != nil {
		slog.Error("error parsing split dividend value ", "string ", splitValues[0])
		return nil, fmt.Errorf("error parsing split dividend value string (%s)", splitValues[0])
	}

	divisor, err := strconv.ParseFloat(splitValues[1], 64)
	if err != nil {
		slog.Error("error parsing split divisor value ", "string ", splitValues[1])
		return nil, fmt.Errorf("error parsing split divisor value string (%s)", splitValues[1])
	}

	if divisor < 0.0000001 {
		slog.Error("split factor divisor cannot be zero (0)", "input", divisor)
		return nil, fmt.Errorf("split factor (%s) divisor cannot be zero (0)", dataPoints[1])
	}

	return &haybaler.EodSplit{
		SplitDate: splitDate,
		Ticker:    ticker,
		Exchange:  exchange,
		Factor:    float32(dividend / divisor),
	}, nil
}
