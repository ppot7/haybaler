package datautils

import "time"

type PriceRangeDataRetriever interface {
	RetrievePriceRangeData(ticker string, exchange string, beginDate time.Time, endDate time.Time, flags ...uint32) ([]PriceRangeData, error)
	RetrieveDividendData(ticker string, exchange string, beginDate time.Time, endDate time.Time, flags ...uint32) ([]DividendData, error)
	RetrieveSplitData(ticker string, exchange string, beginDate time.Time, endDate time.Time, flags ...uint32) ([]SplitData, error)
}

type PriceRangeDataIngestor interface {
	IngestPriceRangeData([]PriceRangeData) error
	IngestDividendData([]DividendData) error
	IngestSplitData([]SplitData) error
}
