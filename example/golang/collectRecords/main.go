package main

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/miaolz123/conver"
	"github.com/miaolz123/stockdb/stockdb"
)

func main() {
	records := [][6]float64{}
	req := httplib.Get("https://www.okcoin.cn/api/v1/kline.do?symbol=btc_cny&type=5min")
	if err := req.ToJSON(&records); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(181818, time.Now().Unix())
		cli := stockdb.New("http://localhost:8765", "username:password")
		opt := stockdb.Option{
			Market: "okcoin_cn",
			Symbol: "BTC_CNY",
			Period: 5 * stockdb.Minute,
		}
		for _, record := range records {
			resp := cli.PutOHLC(stockdb.OHLC{
				Time:   conver.Int64Must(record[0]) / 1000,
				Open:   record[1],
				High:   record[2],
				Low:    record[3],
				Close:  record[4],
				Volume: record[5],
			}, opt)
			if !resp.Success {
				fmt.Println(333333, resp.Message)
			}
		}
		fmt.Println(383838, time.Now().Unix())
	}
}
