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
	return nil
}

func (*EodPsConnection) LoadDividendData(ctx context.Context, dataRange []haybaler.EodDividend) error {
	return nil
}

func (*EodPsConnection) LoadSplitData(ctx context.Context, dataRange []haybaler.EodSplit) error {
	return nil
}
