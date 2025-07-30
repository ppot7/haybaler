package datautils

import (
	"context"
	"time"
)

type PriceRangeDataRetriever interface {
	RetrievePriceRangeData(ticker string, exchange string, beginDate time.Time, endDate time.Time, flags ...uint32) ([]PriceRangeData, error)
	RetrieveDividendData(ticker string, exchange string, beginDate time.Time, endDate time.Time, flags ...uint32) ([]DividendData, error)
	RetrieveSplitData(ticker string, exchange string, beginDate time.Time, endDate time.Time, flags ...uint32) ([]SplitData, error)
}

type PriceRangeDataIngestor interface {
	IngestPriceRangeData(ctx context.Context, priceRange []PriceRangeData) error
	IngestDividendData(ctx context.Context, dividendInfo []DividendData) error
	IngestSplitData(ctx context.Context, splitInfo []SplitData) error
}
