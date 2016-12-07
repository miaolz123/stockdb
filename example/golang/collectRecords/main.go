package main

import (
	"log"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/miaolz123/conver"
	"github.com/miaolz123/stockdb/stockdb"
)

const (
	uri    = "http://localhost:8765"
	auth   = "username:password"
	market = "okcoin_cn"
	symbol = "BTC_CNY"
	period = "1min"
)

func main() {
	periods := map[string]int64{
		"1min": 60,
		"3min": 180,
		"5min": 300,
	}
	if periods[period] == 0 {
		log.Fatalln("period error")
	}
	opt := stockdb.Option{
		Market: market,
		Symbol: symbol,
		Period: periods[period],
	}
	for {
		fetch(opt)
		time.Sleep(200 * time.Duration(opt.Period) * time.Second)
	}
}

func fetch(opt stockdb.Option) {
	records := [][6]float64{}
	req := httplib.Get("https://www.okcoin.cn/api/v1/kline.do?symbol=btc_cny&type=" + period + "&size=1000")
	if err := req.ToJSON(&records); err != nil {
		log.Println("parse json error: ", err)
	} else {
		ohlcs := []stockdb.OHLC{}
		for _, record := range records {
			ohlcs = append(ohlcs, stockdb.OHLC{
				Time:   conver.Int64Must(record[0]) / 1000,
				Open:   record[1],
				High:   record[2],
				Low:    record[3],
				Close:  record[4],
				Volume: record[5],
			})
		}
		if len(ohlcs) > 0 {
			cli := stockdb.New(uri, auth)
			if resp := cli.PutOHLCs(ohlcs, opt); !resp.Success {
				log.Println("PutOHLCs error: ", resp.Message)
			} else {
				log.Println("PutOHLCs successfully")
			}
		}
	}
}
