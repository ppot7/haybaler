package eodps

import (
	"fmt"
	"log/slog"

	"github.com/jackc/pgx"
	"github.com/ppot7/haybaler/datautils"
)

type EodPostgresConnection struct {
	DividendTable string
	PriceTable    string
	SplitsTable   string
	*pgx.Conn
}

func CreateConnectionConfig(host string, port uint16, user string, pwd string, database string) *pgx.ConnConfig {
	return &pgx.ConnConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: pwd,
		Database: database,
	}
}

func CreateEodPostgresConnection(priceTable string, dividendTable string, splitTable string, config *pgx.ConnConfig) (*EodPostgresConnection, error) {
	conn, err := pgx.Connect(*config)
	if err != nil {
		slog.Error("connection error ", "err: ", err)
		return nil, fmt.Errorf("connection error: %s", err)
	}

	return &EodPostgresConnection{
		PriceTable:    priceTable,
		DividendTable: dividendTable,
		SplitsTable:   splitTable,
		Conn:          conn,
	}, nil
}

func (p *EodPostgresConnection) IngestPriceRangeData([]datautils.PriceRangeData) error {
	return nil
}
