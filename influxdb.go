package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/miaolz123/conver"
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

// record2BatchPoints parse struct from OHLC to BatchPoints
func (driver *influxdb) record2BatchPoints(data []stockdb.OHLC, opt stockdb.Option) (bp client.BatchPoints, err error) {
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
	timeOffsets := [4]int64{0, 1, 1, 2}
	for _, datum := range data {
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
			pt, err := client.NewPoint("symbol_"+opt.Symbol, tags[i], fields[i], time.Unix(datum.Time+timeOffsets[i], 0))
			if err != nil {
				return bp, err
			}
			bp.AddPoint(pt)
		}
	}
	return
}

// putMarket create a new market to stockdb
func (driver *influxdb) putMarket(market string) (resp response) {
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
	if opt.Period < minPeriod {
		opt.Period = defaultOption.Period
	}
	bp, err := driver.record2BatchPoints([]stockdb.OHLC{datum}, opt)
	if err != nil {
		log(logError, err)
		resp.Message = err.Error()
		return
	}
	if err := driver.client.Write(bp); err != nil {
		if strings.Contains(err.Error(), "database not found") {
			resp = driver.putMarket(opt.Market)
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
	if opt.Period < minPeriod {
		opt.Period = defaultOption.Period
	}
	bp, err := driver.record2BatchPoints(data, opt)
	if err != nil {
		log(logError, err)
		resp.Message = err.Error()
		return
	}
	if err := driver.client.Write(bp); err != nil {
		if strings.Contains(err.Error(), "database not found") {
			resp = driver.putMarket(opt.Market)
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

// getTimeRange return the first and the last record time
func (driver *influxdb) getTimeRange(opt stockdb.Option) (ranges [2]int64) {
	params := [2]string{"FIRST", "LAST"}
	for i, param := range params {
		raw := fmt.Sprintf(`SELECT %v("price") FROM "symbol_%v"`, param, opt.Symbol)
		q := client.NewQuery(raw, "market_"+opt.Market, "s")
		if response, err := driver.client.Query(q); err == nil && response.Err == "" && len(response.Results) > 0 {
			result := response.Results[0]
			if result.Err == "" && len(result.Series) > 0 && len(result.Series[0].Values) > 0 && len(result.Series[0].Values[0]) > 0 {
				ranges[i] = conver.Int64Must(result.Series[0].Values[0][0])
			}
		}
	}
	return
}

// GetTimeRange return the first and the last record time
func (driver *influxdb) GetTimeRange(opt stockdb.Option) (resp response) {
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
	resp.Data = driver.getTimeRange(opt)
	resp.Success = true
	return
}

// getOHLCQuery return a query of OHLC
func (driver *influxdb) getOHLCQuery(opt stockdb.Option) (q client.Query) {
	ranges := driver.getTimeRange(opt)
	if opt.EndTime <= 0 || opt.EndTime > ranges[1] {
		opt.EndTime = ranges[1]
	}
	if opt.BeginTime <= 0 {
		opt.BeginTime = opt.EndTime - 999*opt.Period
	}
	if opt.BeginTime < ranges[0] {
		opt.BeginTime = ranges[0]
	}
	raw := fmt.Sprintf(`SELECT FIRST("price"), MAX("price"), MIN("price"),
		LAST("price"), SUM("amount") FROM "symbol_%v" WHERE "period" <= %v
		AND time >= %vs AND time < %vs GROUP BY time(%vs)`, opt.Symbol,
		opt.Period, opt.BeginTime, opt.EndTime, opt.Period)
	q = client.NewQuery(raw, "market_"+opt.Market, "s")
	return q
}

// recordResult2ohlc parse record result to OHLC
func (driver *influxdb) recordResult2ohlc(result client.Result, opt stockdb.Option) (data []stockdb.OHLC) {
	if len(result.Series) > 0 {
		serie := result.Series[0]
		for i := range serie.Values {
			if conver.Float64Must(serie.Values[i][3]) <= 0.0 {
				if opt.InvalidPolicy != "ibid" {
					continue
				} else {
					i--
				}
				if i < 0 {
					continue
				}
			}
			data = append(data, stockdb.OHLC{
				Time:   conver.Int64Must(serie.Values[i][0]),
				Open:   conver.Float64Must(serie.Values[i][1]),
				High:   conver.Float64Must(serie.Values[i][2]),
				Low:    conver.Float64Must(serie.Values[i][3]),
				Close:  conver.Float64Must(serie.Values[i][4]),
				Volume: conver.Float64Must(serie.Values[i][5]),
			})
		}
	}
	return
}

// GetOHLC get OHLC records
func (driver *influxdb) GetOHLCs(opt stockdb.Option) (resp response) {
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
	if opt.Period < minPeriod {
		opt.Period = defaultOption.Period
	}
	if response, err := driver.client.Query(driver.getOHLCQuery(opt)); err != nil {
		log(logError, err)
		resp.Message = err.Error()
		return
	} else if response.Err != "" {
		log(logError, response.Err)
		resp.Message = response.Err
		return
	} else if len(response.Results) > 0 {
		result := response.Results[0]
		if result.Err != "" {
			log(logError, response.Err)
			resp.Message = response.Err
			return
		}
		resp.Data = driver.recordResult2ohlc(result, opt)
		resp.Success = true
	}
	return
}

/*
SELECT FIRST("price") AS open, MAX("price") AS high, MIN("price") AS low, LAST("price") AS close, SUM("amount") AS volume FROM "symbol_BTC_CNY" WHERE time >= '2016-11-23T14:00:00Z' AND time <= '2016-11-27T14:00:00Z' GROUP BY time(30m)
*/
