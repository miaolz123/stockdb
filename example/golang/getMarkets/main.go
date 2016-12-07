package main

import (
	"fmt"

	"github.com/miaolz123/stockdb/stockdb"
)

func main() {
	cli := stockdb.New("http://localhost:8765", "username:password")
	resp := cli.GetMarkets()
	if len(resp.Data) > 0 {
		fmt.Printf("Markets: %+v\n", resp.Data)
		fmt.Printf("Symbols of %s: %+v\n", resp.Data[0], cli.GetSymbols(resp.Data[0]).Data)
	}
}
