package main

import (
	"flag"
	"fmt"

	"github.com/hprose/hprose-golang/io"
	"github.com/miaolz123/stockdb/stockdb"
)

const (
	version = "0.0.1"
)

func init() {
	io.Register(response{}, "response", "json")
	io.Register(stockdb.Option{}, "option", "json")
	io.Register(stockdb.Ticker{}, "ticker", "json")
	io.Register(stockdb.OHLC{}, "ohlc", "json")
	confPath := flag.String("conf", "stockdb.ini", "config file path")
	flag.Parse()
	loadConfig(*confPath)
	log(logInfo, fmt.Sprintf("StockDB Version %s running at %s", version, config["http.bind"]))
}
