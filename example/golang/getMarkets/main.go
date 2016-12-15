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
		fmt.Printf("Symbols of %s: %+v\n", market, cli.GetSymbols(market).Data)
	}
}
