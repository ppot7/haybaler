package haybaler

import (
	"context"
	"iter"
	"time"
)

type EodDataRetriever interface {
	RetrievePriceVolumeData(ctx context.Context, ticker string, exchange string, begin time.Time, end time.Time, options ...uint32) (iter.Seq2[*EodPriceVolume, error], error)
	RetrieveDividendData(ctx context.Context, ticker string, exchange string, begin time.Time, end time.Time, options ...uint32) (iter.Seq2[*EodDividend, error], error)
	RetrieveSplitData(ctx context.Context, ticker string, exchange string, begin time.Time, end time.Time, optons ...uint32) (iter.Seq2[*EodSplit, error], error)
}

type EodDataLoader interface {
	LoadPriceVolumeData(ctx context.Context, dataRange []EodPriceVolume) error
	LoadDividendData(ctx context.Context, dataRange []EodDividend) error
	LoadSplitData(ctx context.Context, dataRange []EodSplit) error
}
