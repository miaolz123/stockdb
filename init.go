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
	io.Register(response{}, "Response", "json")
	io.Register(stockdb.Option{}, "Option", "json")
	io.Register(stockdb.Ticker{}, "Ticker", "json")
	io.Register(stockdb.OHLC{}, "OHLC", "json")
	io.Register(stockdb.OrderBook{}, "OrderBook", "json")
	io.Register(stockdb.Depth{}, "Depth", "json")
	confPath := flag.String("conf", "stockdb.ini", "config file path")
	flag.Parse()
	loadConfig(*confPath)
	log(logInfo, fmt.Sprintf("StockDB Version %s running at %s", version, config["http.bind"]))
}
