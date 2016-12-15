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
	uri    = "http://localhost:8765"
	auth   = "username:password"
	market = "haobtc"
	symbol = "150W_HASH/BTC"
)

var location *time.Location

func main() {
	if loc, err := time.LoadLocation("Asia/Shanghai"); err != nil || loc == nil {
		location = time.Local
	} else {
		location = loc
	}
	opt := stockdb.Option{
		Market: market,
		Symbol: symbol,
	}
	for {
		fetch(opt)
		time.Sleep(10 * time.Minute)
	}
}

func fetch(opt stockdb.Option) {
	req := httplib.Get("https://hashex.haobtc.com/exchange/main?code=200001&limit=undefined")
	req.Header("Cookie", "")
	if resp, err := req.Bytes(); err != nil {
		log.Println("http error: ", err)
	} else {
		if json, err := simplejson.NewJson(resp); err != nil {
			log.Println("parse json error: ", err)
		} else {
			if json.Get("result").MustString() != "success" {
				log.Println("get data error: ", json)
			}
			records := json.GetPath("data", "public", "trades")
			orders := []stockdb.Order{}
			for i := 0; i < len(records.MustArray()); i++ {
				record := records.GetIndex(i).MustMap()
				t, err := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprint(record["time"]), location)
				if err != nil {
					log.Println("parse time error: ", err)
					continue
				}
				orders = append(orders, stockdb.Order{
					ID:     fmt.Sprint(record["volume"], "@", conver.Float64Must(record["price"])),
					Time:   t.Unix(),
					Price:  conver.Float64Must(record["price"]),
					Amount: conver.Float64Must(record["volume"]),
					Type:   fmt.Sprint(record["side"]),
				})
			}
			if len(orders) > 0 {
				cli := stockdb.New(uri, auth)
				if resp := cli.PutOrders(orders, opt); !resp.Success {
					log.Println("PutOrders error: ", resp.Message)
				} else {
					log.Println("PutOrders successfully")
				}
			}
		}
	}
}
