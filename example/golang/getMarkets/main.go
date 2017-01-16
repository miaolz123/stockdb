package main

import (
	"fmt"

	"github.com/miaolz123/stockdb/stockdb"
)

func main() {
	cli := stockdb.New("http://localhost:8765", "username:password")
	fmt.Printf("%+v\n", cli.GetStats())
	resp := cli.GetMarkets()
	for _, market := range resp.Data {
		symbols := cli.GetSymbols(market).Data
		fmt.Printf("Symbols of %s: %+v\n", market, symbols)
		for _, symbol := range symbols {
			fmt.Printf("MinPeriod of %s: %+v\n", symbol, cli.GetPeriodRange(stockdb.Option{
				Market: market,
				Symbol: symbol,
			}).Data)
		}
	}
}
