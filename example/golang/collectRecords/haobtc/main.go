package main

import (
	"fmt"
	"log"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/miaolz123/stockdb/stockdb"
)

const (
	uri    = "http://localhost:8765"
	auth   = "username:password"
	market = "haobtc"
	symbol = "BTC/CNY"
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
	records := []struct {
		Volume float64
		Price  float64
		Side   string
		Time   string
	}{}
	req := httplib.Get("https://haobtc.com/exchange/api/v1/trades")
	if err := req.ToJSON(&records); err != nil {
		log.Println("parse json error: ", err)
	} else {
		orders := []stockdb.Order{}
		for _, record := range records {
			t, err := time.ParseInLocation("2006-01-02 15:04:05", record.Time, location)
			if err != nil {
				log.Println("parse time error: ", err)
				continue
			}
			orders = append(orders, stockdb.Order{
				ID:     fmt.Sprint(record.Volume, "@", record.Price),
				Time:   t.Unix(),
				Price:  record.Price,
				Amount: record.Volume,
				Type:   record.Side,
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
