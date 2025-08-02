package eodpostgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type EodPsConnection struct {
	pgx.Conn
	SchemaName       string
	PriceVolumeTable string
	DividendTable    string
	SplitTable       string
}

func CreateDefaultConfiguration(host string, port uint, user string, password string, database string) *pgx.ConnConfig {
	config := &pgconn.Config{
		Host:     host,
		Port:     uint16(port),
		User:     user,
		Password: password,
		Database: database,
	}

	return &pgx.ConnConfig{
		Config: *config,
	}
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
	}, nil
}

func (*EodPsConnection) CreateSchema() error {
	return nil
}

func (*EodPsConnection) DropSchema() error {
	return nil
}

func (*EodPsConnection) CreatePriceVolumeTable() error {
	return nil
}

func (*EodPsConnection) DropPriceVolumeTable() error {
	return nil
}

func (*EodPsConnection) CreateDividendTable() error {
	return nil
}

func (*EodPsConnection) DropDividendTable() error {
	return nil
}

func (*EodPsConnection) CreateSplitTable() error {
	return nil
}

func (*EodPsConnection) DropSplitTable() error {
	return nil
}
