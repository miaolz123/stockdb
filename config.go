package main

import (
	"encoding/base64"
	"strings"

	"github.com/go-ini/ini"
)

var config = make(map[string]string)

func loadConfig(path string) {
	if path == "" {
		path = "default.ini"
	}
	conf, err := ini.Load(path)
	if err != nil {
		log(logFatal, "Load config file error: ", err)
	}
	for _, s := range conf.Sections() {
		name := s.Name()
		for _, k := range s.Keys() {
			config[strings.ToLower(name+"."+k.Name())] = k.Value()
		}
	}
	config["http.auth"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(config["http.auth"]))
}
