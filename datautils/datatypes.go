package datautils

import (
	"fmt"
	"time"
)

const (
	PRICE    uint32 = (1 << 0)
	DIVIDEND uint32 = (1 << 1)
	SPLIT    uint32 = (1 << 2)
	CSV      uint32 = (1 << 3)
	JSON     uint32 = (1 << 5)
)

type PriceRangeData struct {
	TradeDate time.Time
	Ticker    string
	Exchange  string
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    int64
}

type DividendData struct {
	ExDate        time.Time
	Ticker        string
	Exchange      string
	DividendValue float64
}

type SplitData struct {
	SplitDate   time.Time
	Ticker      string
	Exchange    string
	SplitFactor float64
}

func (p *PriceRangeData) GoString() string {
	dataStr := fmt.Sprintf("Date: %s  Ticker: %s  Exchange: %s  Open: %.5f  High: %.5f  Low: %.5f  Close: %.5f  Volume:  %d",
		p.TradeDate.Format(time.DateOnly), p.Ticker, p.Exchange, p.Open, p.High, p.Low, p.Close, p.Volume)
	return dataStr
}

func (p *DividendData) GoString() string {
	return fmt.Sprintf("Ex-Date: %s  Ticker:  %s  Exchange: %s  Dividend: %.5f", p.ExDate.Format(time.DateOnly), p.Ticker, p.ExDate, p.DividendValue)
}

func (p *SplitData) GoString() string {
	return fmt.Sprintf("Split Date: %s  Ticker: %s  Exchange: %s  Factor: %.5f", p.SplitDate.Format(time.DateOnly), p.Ticker, p.Exchange, p.SplitFactor)
}
