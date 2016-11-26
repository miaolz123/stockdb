package main

import (
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

type influxdb struct {
	client client.Client
	status int64
}

// newInfluxdb create a Influxdb struct
func newInfluxdb() driver {
	var err error
	driver := &influxdb{}
	driver.client, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     config["influxdb.host"],
		Username: config["influxdb.username"],
		Password: config["influxdb.password"],
	})
	if err != nil {
		log(logFatal, "Influxdb connect error: ", err)
	}
	if _, _, err := driver.client.Ping(30 * time.Second); err != nil {
		log(logError, "Influxdb connect error: ", err)
	} else {
		driver.status = 1
		log(logSuccess, "Influxdb connect successfully")
	}
	go func(driver *influxdb) {
		for {
			if _, _, err := driver.client.Ping(30 * time.Second); err != nil {
				driver.status = 0
				driver.reconnect()
			}
			time.Sleep(time.Minute)
		}
	}(driver)
	return driver
}

// reconnect the client
func (driver *influxdb) reconnect() {
	for {
		time.Sleep(10 * time.Minute)
		var err error
		if driver.client, err = client.NewHTTPClient(client.HTTPConfig{
			Addr:     config["influxdb.host"],
			Username: config["influxdb.username"],
			Password: config["influxdb.password"],
		}); err == nil {
			if _, _, err = driver.client.Ping(30 * time.Second); err == nil {
				log(logSuccess, "Influxdb reconnect successfully")
				driver.status = 1
				break
			}
		}
		log(logError, "Influxdb reconnect error: ", err)
	}
}

// close this client
func (driver *influxdb) close() error {
	log(logSuccess, "Influxdb disconnected successfully")
	return driver.client.Close()
}

// check the client is connected
func (driver *influxdb) check() error {
	if driver.status < 1 {
		return errInfluxdbNotConnected
	}
	return nil
}

// ohlc2BatchPoints parse struct from OHLC to BatchPoints
func (driver *influxdb) ohlc2BatchPoints(market, symbol string, data []ohlc) (bp client.BatchPoints, err error) {
	if driver.status < 1 {
		err = errInfluxdbNotConnected
		return
	}
	bp, err = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  market,
		Precision: "s",
	})
	if err != nil {
		return
	}
	for _, datum := range data {
		time := time.Unix(datum.Time, 0)
		tags := [4]map[string]string{
			{"type": "open"},
			{"type": "high"},
			{"type": "low"},
			{"type": "close"},
		}
		fields := [4]map[string]interface{}{
			{
				"price":  datum.Open,
				"volume": datum.Volume / 4.0,
			},
			{
				"price":  datum.High,
				"volume": datum.Volume / 4.0,
			},
			{
				"price":  datum.Low,
				"volume": datum.Volume / 4.0,
			},
			{
				"price":  datum.Close,
				"volume": datum.Volume / 4.0,
			},
		}
		for i := 0; i < 4; i++ {
			pt, err := client.NewPoint(symbol, tags[i], fields[i], time)
			if err != nil {
				return bp, err
			}
			bp.AddPoint(pt)
		}
	}
	return
}

// AddMarket create a new market to stockdb
func (driver *influxdb) AddMarket(market string) (resp response) {
	if err := driver.check(); err != nil {
		log(logError, err)
		resp.Message = err.Error()
		return
	}
	q := client.NewQuery("CREATE DATABASE "+market, "", "")
	if response, err := driver.client.Query(q); err != nil {
		log(logError, err)
		resp.Message = err.Error()
		return
	} else if err = response.Error(); err != nil {
		log(logError, err)
		resp.Message = err.Error()
		return
	}
	resp.Success = true
	return
}

// AddOHLC add a OHLC record to stockdb
func (driver *influxdb) AddOHLC(market, symbol string, datum ohlc) (resp response) {
	if err := driver.check(); err != nil {
		log(logError, err)
		resp.Message = err.Error()
		return
	}
	bp, err := driver.ohlc2BatchPoints(market, symbol, []ohlc{datum})
	if err != nil {
		log(logError, err)
		resp.Message = err.Error()
		return
	}
	if err := driver.client.Write(bp); err != nil {
		if strings.Contains(err.Error(), "database not found") {
			resp = driver.AddMarket(market)
			if resp.Success {
				return driver.AddOHLC(market, symbol, datum)
			}
			return
		}
		log(logError, err)
		resp.Message = err.Error()
		return
	}
	resp.Success = true
	return
}

// AddOHLC add a OHLC record to stockdb
func (driver *influxdb) AddOHLCs(market, symbol string, data []ohlc) (resp response) {
	if err := driver.check(); err != nil {
		log(logError, err)
		resp.Message = err.Error()
		return
	}
	bp, err := driver.ohlc2BatchPoints(market, symbol, data)
	if err != nil {
		log(logError, err)
		resp.Message = err.Error()
		return
	}
	if err := driver.client.Write(bp); err != nil {
		if strings.Contains(err.Error(), "database not found") {
			resp = driver.AddMarket(market)
			if resp.Success {
				return driver.AddOHLCs(market, symbol, data)
			}
			return
		}
		log(logError, err)
		resp.Message = err.Error()
		return
	}
	resp.Success = true
	return
}

/*
SELECT FIRST("price") AS open, MAX("price") AS high, MIN("price") AS low, LAST("price") AS close, SUM("volume") AS volume FROM "symbol" WHERE price > 0 AND time >= '2016-11-26T14:00:00Z' GROUP BY time(3m)
*/
