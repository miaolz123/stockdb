package main

import (
	"flag"
	"fmt"

	"github.com/hprose/hprose-golang/io"
)

const (
	version         = "0.1.4"
	minPeriod int64 = 3
)

func init() {
	io.Register(response{}, "Response", "json")
	confPath := flag.String("conf", "stockdb.ini", "config file path")
	flag.Parse()
	loadConfig(*confPath)
	log(logInfo, fmt.Sprintf("StockDB Version %s running at %s", version, config["http.bind"]))
}
