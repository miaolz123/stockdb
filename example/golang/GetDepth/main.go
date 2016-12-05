package main

import (
	"fmt"

	"github.com/miaolz123/stockdb/stockdb"
)

func main() {
	cli := stockdb.New("http://localhost:8765", "username:password")
	opt := stockdb.Option{
		BeginTime: 1479916800,
		Period:    stockdb.Minute * 30,
	}
	fmt.Printf("%+v\n", cli.GetDepth(opt))
}
