package haybaler

import (
	"context"
	"time"
)

type EodDataRetriever interface {
	RetrievePriceVolumeData(ctx context.Context, ticker string, exchange string, begin time.Time, end time.Time, options ...uint32) ([]EodPriceVolume, error)
	RetrieveDividendData(ctx context.Context, ticker string, exchange string, begin time.Time, end time.Time, options ...uint32) ([]EodDividend, error)
	RetrieveSplitData(ctx context.Context, ticker string, exchange string, begin time.Time, end time.Time, optons ...uint32) ([]EodSplit, error)
}

type EodDataLoader interface {
	LoadPriceVolumeData(ctx context.Context, dataRange []EodPriceVolume) error
	LoadDividendData(ctx context.Context, dataRange []EodDividend) error
	LoadSplitData(ctx context.Context, dataRange []EodSplit) error
}
