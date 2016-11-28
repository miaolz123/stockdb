package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/miaolz123/stockdb/stockdb"
)

type influxdb struct {
	client client.Client
	status int64
}

// newInfluxdb create a Influxdb struct
func newInfluxdb() Driver {
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
func (driver *influxdb) ohlc2BatchPoints(data []stockdb.OHLC, opt stockdb.Option) (bp client.BatchPoints, err error) {
	if driver.status < 1 {
		err = errInfluxdbNotConnected
		return
	}
	bp, err = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "market_" + opt.Market,
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
			{"price": datum.Open},
			{"price": datum.High},
			{"price": datum.Low},
			{"price": datum.Close},
		}
		for i := 0; i < 4; i++ {
			tags[i]["id"] = fmt.Sprint(opt.Period)
			fields[i]["period"] = opt.Period
			fields[i]["amount"] = datum.Volume / 4.0
			pt, err := client.NewPoint("symbol_"+opt.Symbol, tags[i], fields[i], time)
			if err != nil {
				return bp, err
			}
			bp.AddPoint(pt)
		}
	}
	return
}

// PutMarket create a new market to stockdb
func (driver *influxdb) PutMarket(market string) (resp response) {
	if err := driver.check(); err != nil {
		log(logError, err)
		resp.Message = err.Error()
		return
	}
	q := client.NewQuery("CREATE DATABASE market_"+market, "", "")
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

// PutOHLC add a OHLC record to stockdb
func (driver *influxdb) PutOHLC(datum stockdb.OHLC, opt stockdb.Option) (resp response) {
	if err := driver.check(); err != nil {
		log(logError, err)
		resp.Message = err.Error()
		return
	}
	if opt.Market == "" {
		opt.Market = defaultOption.Market
	}
	if opt.Symbol == "" {
		opt.Symbol = defaultOption.Symbol
	}
	if opt.Period == 0 {
		opt.Period = defaultOption.Period
	}
	bp, err := driver.ohlc2BatchPoints([]stockdb.OHLC{datum}, opt)
	if err != nil {
		log(logError, err)
		resp.Message = err.Error()
		return
	}
	if err := driver.client.Write(bp); err != nil {
		if strings.Contains(err.Error(), "database not found") {
			resp = driver.PutMarket(opt.Market)
			if resp.Success {
				return driver.PutOHLC(datum, opt)
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

// PutOHLC add a OHLC record to stockdb
func (driver *influxdb) PutOHLCs(data []stockdb.OHLC, opt stockdb.Option) (resp response) {
	if err := driver.check(); err != nil {
		log(logError, err)
		resp.Message = err.Error()
		return
	}
	if opt.Market == "" {
		opt.Market = defaultOption.Market
	}
	if opt.Symbol == "" {
		opt.Symbol = defaultOption.Symbol
	}
	if opt.Period == 0 {
		opt.Period = defaultOption.Period
	}
	bp, err := driver.ohlc2BatchPoints(data, opt)
	if err != nil {
		log(logError, err)
		resp.Message = err.Error()
		return
	}
	if err := driver.client.Write(bp); err != nil {
		if strings.Contains(err.Error(), "database not found") {
			resp = driver.PutMarket(opt.Market)
			if resp.Success {
				return driver.PutOHLCs(data, opt)
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
SELECT FIRST("price") AS open, MAX("price") AS high, MIN("price") AS low, LAST("price") AS close, SUM("amount") AS volume FROM "symbol_BTC_CNY" WHERE time >= '2016-11-23T14:00:00Z' AND time <= '2016-11-27T14:00:00Z' GROUP BY time(30m)
*/
