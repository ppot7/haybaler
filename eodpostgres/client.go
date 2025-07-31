package eodpostgres

import "github.com/jackc/pgx"

type EodPostgresConnection struct {
	pgx.Conn
	Schema      string
	PriceVolume string
	Dividends   string
	Splits      string
}
