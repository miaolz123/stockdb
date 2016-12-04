package main

import (
	"encoding/base64"
	"strings"

	"github.com/go-ini/ini"
	"github.com/miaolz123/stockdb/stockdb"
)

var (
	config        = make(map[string]string)
	openMethods   = make(map[string]bool)
	defaultOption = stockdb.Option{}
)

func loadConfig(path string) {
	if path == "" {
		path = "stockdb.ini"
	}
	conf, err := ini.Load(path)
	if err != nil {
		log(logFatal, "Load config file error: ", err)
	}
	conf.Section("default").MapTo(&defaultOption)
	if defaultOption.Period < minPeriod {
		defaultOption.Period = minPeriod
	}
	for _, s := range conf.Sections() {
		name := s.Name()
		for _, k := range s.Keys() {
			config[strings.ToLower(name+"."+k.Name())] = k.Value()
		}
	}
	for _, m := range strings.Split(config["http.openmethods"], ",") {
		openMethods[strings.TrimSpace(m)] = true
	}
	config["http.auth"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(config["http.auth"]))
}
