package eodps

var priceRangeTableBase string = `CREATE TABLE IF NOT EXISTS "%s"."%s"
(
    "TRADE_DATE" date NOT NULL,
    "SYMBOL" character varying COLLATE pg_catalog."default" NOT NULL,
    "EXCHANGE" character varying COLLATE pg_catalog."default" NOT NULL,
    "OPEN" numeric,
    "HIGH" numeric,
    "LOW" numeric,
    "CLOSE" numeric,
	"VOLUME" integer,
    "_CREATE_AT" timestamp with time zone,
    "_UPDATE_AT" timestamp with time zone,
    CONSTRAINT "DAILY_PRICE_RANGE_VOL_pkey" PRIMARY KEY ("TRADE_DATE", "SYMBOL", "EXCHANGE")
)
TABLESPACE pg_default;
`
