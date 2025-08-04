package haybaler

import (
	"fmt"
	"time"
)

func (p *EodPriceVolume) GoString() string {
	dataStr := fmt.Sprintf("Date: %s  Ticker: %s  Exchange: %s  Open: %.5f  High: %.5f  Low: %.5f  Close: %.5f  Volume:  %d",
		p.TradeDate.Format(time.DateOnly), p.Ticker, p.Exchange, p.Open, p.High, p.Low, p.Close, p.Volume)
	return dataStr
}

func (p *EodDividend) GoString() string {
	return fmt.Sprintf("Ex-Date: %s  Ticker:  %s  Exchange: %s  Dividend: %.5f", p.ExDate.Format(time.DateOnly), p.Ticker, p.Exchange, p.Value)
}

func (p *EodSplit) GoString() string {
	return fmt.Sprintf("Split Date: %s  Ticker: %s  Exchange: %s  Factor: %.5f", p.SplitDate.Format(time.DateOnly), p.Ticker, p.Exchange, p.Factor)
}
