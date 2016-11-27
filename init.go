package main

import (
	"flag"
	"fmt"

	"github.com/hprose/hprose-golang/io"
)

const (
	version = "0.0.1"
)

func init() {
	io.Register(response{}, "response", "json")
	io.Register(option{}, "option", "json")
	io.Register(ticker{}, "ticker", "json")
	io.Register(ohlc{}, "ohlc", "json")
	confPath := flag.String("conf", "default.ini", "config file path")
	flag.Parse()
	loadConfig(*confPath)
	log(logInfo, fmt.Sprintf("StockDB Version %s running at %s", version, config["http.bind"]))
}
