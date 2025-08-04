package eodpostgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/ppot7/haybaler"
)

type EodPsConnection struct {
	pgx.Conn
	SchemaName       string
	PriceVolumeTable string
	DividendTable    string
	SplitTable       string
}

func CreateDefaultConfiguration(host string, port string, user string, password string, database string) (*pgx.ConnConfig, error) {
	/* ParseConfig must be used to create a configuration or panic is thrown. */
	config, err := pgx.ParseConfig(fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s",
		host, user, password, port, database))
	if err != nil {
		slog.Error("could not parse connection string", "err", err)
		return nil, fmt.Errorf("configuration error: %s", err)
	}

	return config, nil
}

func ConnectToPsDatabase(ctx context.Context, schema string, priceVol string, dividend string, split string, config *pgx.ConnConfig) (*EodPsConnection, error) {
	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		slog.Error("unable to establish connection", "err", err)
		return nil, fmt.Errorf("unable to establish connection: %s", err)
	}

	return &EodPsConnection{
		Conn:             *conn,
		SchemaName:       schema,
		PriceVolumeTable: priceVol,
		DividendTable:    dividend,
		SplitTable:       split,
	}, nil
}

func (p *EodPsConnection) LoadPriceVolumeData(ctx context.Context, dataRange []haybaler.EodPriceVolume) error {

	query_base := fmt.Sprintf(`INSERT INTO "%s"."%s " `, p.SchemaName, p.DividendTable)
	query_base += `(trade_date, ticker, exchange, open_price, high_price, low_price, close_price, volume)`
	query_base += ` VALUES (@date, @ticker, @exchange, @open, @high, @low, @close, @volume)`

	batch := &pgx.Batch{}
	for _, data := range dataRange {

		batch.Queue(query_base, pgx.NamedArgs{
			"date":     data.TradeDate,
			"ticker":   data.Ticker,
			"exchange": data.Exchange,
			"open":     data.Open,
			"high":     data.High,
			"low":      data.Low,
			"close":    data.Close,
			"volume":   data.Volume,
		})
	}

	batchResults := p.SendBatch(ctx, batch)
	defer batchResults.Close()

	errCount := 0
	for i := range len(dataRange) {
		_, err := batchResults.Exec()
		if err != nil {
			errCount++
			slog.Error("batch insert error", "err", err, "record", dataRange[i].GoString())
		}
	}

	if errCount > 0 {
		return fmt.Errorf("%d non-fatal batch insert errors found (see logs)", errCount)
	}

	return nil
}

func (p *EodPsConnection) LoadDividendData(ctx context.Context, dataRange []haybaler.EodDividend) error {

	query_base := fmt.Sprintf(`INSERT INTO "%s"."%s" (ex_dividend_date, ticker, exchange, dividend_value)`, p.SchemaName, p.DividendTable)
	query_base += ` VALUES (@date, @ticker, @exchange, @dividend)`

	batch := &pgx.Batch{}
	for _, data := range dataRange {

		batch.Queue(query_base, pgx.NamedArgs{
			"date":     data.ExDate,
			"ticker":   data.Ticker,
			"exchange": data.Exchange,
			"factor":   data.Value,
		})
	}

	batchResults := p.SendBatch(ctx, batch)
	defer batchResults.Close()

	errCount := 0
	for i := range len(dataRange) {
		_, err := batchResults.Exec()
		if err != nil {
			errCount++
			slog.Error("batch insert error", "err", err, "record", dataRange[i].GoString())
		}
	}

	if errCount > 0 {
		return fmt.Errorf("%d non-fatal batch insert errors found (see logs)", errCount)
	}

	return nil
}

func (p *EodPsConnection) LoadSplitData(ctx context.Context, dataRange []haybaler.EodSplit) error {

	query_base := fmt.Sprintf(`INSERT INTO "%s"."%s" (split_date, ticker, exchange, split_factor)`, p.SchemaName, p.SplitTable)
	query_base += ` VALUES (@date, @ticker, @exchange, @factor)`

	batch := &pgx.Batch{}
	for _, data := range dataRange {

		batch.Queue(query_base, pgx.NamedArgs{
			"date":     data.SplitDate,
			"ticker":   data.Ticker,
			"exchange": data.Exchange,
			"factor":   data.Factor,
		})
	}

	batchResults := p.SendBatch(ctx, batch)
	defer batchResults.Close()

	errCount := 0
	for i := range len(dataRange) {
		_, err := batchResults.Exec()
		if err != nil {
			errCount++
			slog.Error("batch insert error", "err", err, "record", dataRange[i].GoString())
		}
	}

	if errCount > 0 {
		return fmt.Errorf("%d non-fatal batch insert errors found (see logs)", errCount)
	}

	return nil
}
