package main

import (
	"flag"
	"fmt"
)

const (
	version = "0.0.1"
)

func init() {
	confPath := flag.String("conf", "default.ini", "config file path")
	flag.Parse()
	loadConfig(*confPath)
	log(logInfo, fmt.Sprintf("StockDB Version %s running at %s", version, config["http.bind"]))
}
