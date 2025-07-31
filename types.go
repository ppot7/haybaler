package haybaler

import (
	"net/http"
	"time"
)

type EodHttpClient struct {
	http.Client
	Url string
}

type EodPriceVolume struct {
	TradeDate time.Time
	Ticker    string
	Exchange  string
	Open      float32
	High      float32
	Low       float32
	Close     float32
	Volume    int32
}

type EodDividend struct {
	ExDate   time.Time
	Ticker   string
	Exchange string
	Value    float32
}

type EodSplit struct {
	SplitDate time.Time
	Ticker    string
	Exchange  string
	Factor    float32
}
