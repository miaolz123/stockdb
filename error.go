package main

import (
	"fmt"
)

var (
	errHTTPUnauthorized     = fmt.Errorf("Unauthorized")
	errInfluxdbNotConnected = fmt.Errorf("Influxdb is not connected")
)
