package eodps

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx"
	"github.com/ppot7/haybaler/datautils"
)

type EodDatabaseSchema struct {
	Schema          string
	PriceRangeTable string
	DividendTable   string
	SplitTable      string
}

type EodPostgresConnection struct {
	Schema EodDatabaseSchema
	*pgx.Conn
}

var defaultEodSchema = EodDatabaseSchema{
	Schema:          "TEST_DAILY_DATA",
	PriceRangeTable: "DAILY_PRICE_RANGE",
	DividendTable:   "DAILY_DIVIDEND",
	SplitTable:      "DAILY_SPLIT",
}

func CreateDefaultEodConnection(host string, port uint16, user string, pwd string, database string) (*EodPostgresConnection, error) {
	return CreateCustomEodConnection(host, port, user, pwd, database, defaultEodSchema)
}

func CreateCustomEodConnection(host string, port uint16, user string, pwd string, database string, eodSchema EodDatabaseSchema) (*EodPostgresConnection, error) {
	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: pwd,
		Database: database,
	})
	if err != nil {
		slog.Error("error connecting to EodConnection", "err", err)
		return nil, err
	}

	return &EodPostgresConnection{
		Schema: eodSchema,
		Conn:   conn,
	}, nil
}

func (p *EodPostgresConnection) IngestPriceRangeData(ctx context.Context, priceRange []datautils.PriceRangeData) error {

	queryBase := fmt.Sprintf(`INSERT INTO "%s"."%s"("TRADE_DATE", "SYMBOL", "EXCHANGE", `, p.Schema.Schema, p.Schema.PriceRangeTable)
	queryBase += `"OPEN", "HIGH", "LOW", "CLOSE", "VOLUME", "_CREATE_AT", "_UPDATE_AT") VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	data := make([]interface{}, 10)
	data[0] = priceRange[0].TradeDate
	data[1] = priceRange[0].Ticker
	data[2] = priceRange[0].Exchange
	data[3] = priceRange[0].Open
	data[4] = priceRange[0].High
	data[5] = priceRange[0].Low
	data[6] = priceRange[0].Close
	data[7] = priceRange[0].Volume
	data[8] = time.Now()
	data[9] = time.Now()

	commTag, err := p.Exec(queryBase, priceRange[0].TradeDate, priceRange[0].Ticker, priceRange[0].Exchange, priceRange[0].Open,
		priceRange[0].High, priceRange[0].Low, priceRange[0].Close, priceRange[0].Volume, time.Now(), time.Now())
	if err != nil {
		slog.Error("error inserting record", "err", err)
		return err
	}

	fmt.Println(commTag)

	return nil
}

func (p *EodPostgresConnection) CreatePriceRangeTable(schema string, tableName string) error {
	_, err := p.Exec(fmt.Sprintf(priceRangeTableBase, schema, tableName))
	if err != nil {
		slog.Error("error creating schema.table ", "err", err)
	}

	return err
}
