package main

import (
	"fmt"

	"github.com/miaolz123/stockdb/stockdb"
)

func main() {
	cli := stockdb.New("http://localhost:8765", "username:password")
	opt := stockdb.Option{Period: stockdb.Hour}
	fmt.Printf("%+v\n", cli.GetTimeRange(opt))
	fmt.Printf("%+v\n", cli.GetOHLCs(opt))
}
