package main

import (
	"fmt"
	"log"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/bitly/go-simplejson"
	"github.com/miaolz123/conver"
	"github.com/miaolz123/stockdb/stockdb"
)

const (
	apiKey = "XXXXXXXXXXXXXXXXXXXXXXXXX" // www.quandl.com
	uri    = "http://localhost:8765"
	auth   = "username:password"
	market = "NASDAQ"
	symbol = "AAPL"
)

var location *time.Location

func main() {
	if loc, err := time.LoadLocation("America/New_York"); err != nil || loc == nil {
		location = time.Local
	} else {
		location = loc
	}
	opt := stockdb.Option{
		Market: market,
		Symbol: symbol,
		Period: stockdb.Day,
	}
	fetch(opt)
}

func fetch(opt stockdb.Option) {
	req := httplib.Get(fmt.Sprintf("https://www.quandl.com/api/v3/datasets/WIKI/%v.json?api_key=%v", symbol, apiKey))
	if resp, err := req.Bytes(); err != nil {
		log.Println("http error: ", err)
	} else {
		if json, err := simplejson.NewJson(resp); err != nil {
			log.Println("parse json error: ", err)
		} else {
			records := json.GetPath("dataset", "data")
			ohlcs := []stockdb.OHLC{}
			for i := 0; i < len(records.MustArray()); i++ {
				record := records.GetIndex(i).MustArray()
				t, err := time.ParseInLocation("2006-01-02", fmt.Sprint(record[0]), location)
				if err != nil {
					log.Println("parse time error: ", err)
					continue
				}
				ohlcs = append(ohlcs, stockdb.OHLC{
					Time:   t.Unix(),
					Open:   conver.Float64Must(record[1]),
					High:   conver.Float64Must(record[2]),
					Low:    conver.Float64Must(record[3]),
					Close:  conver.Float64Must(record[4]),
					Volume: conver.Float64Must(record[5]),
				})
			}
			if len(ohlcs) > 0 {
				queue := len(ohlcs) / 500
				for i := 0; i <= queue; i++ {
					time.Sleep(10 * time.Second)
					begin := 500 * i
					end := begin + 500
					if end > len(ohlcs) {
						end = len(ohlcs)
					}
					cli := stockdb.New(uri, auth)
					if resp := cli.PutOHLCs(ohlcs[begin:end], opt); !resp.Success {
						log.Println("PutOHLCs error: ", resp.Message)
					} else {
						log.Println("PutOHLCs successfully")
					}
				}
			}
		}
	}
}
