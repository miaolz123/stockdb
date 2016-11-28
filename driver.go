package main

import (
	"github.com/miaolz123/stockdb/stockdb"
)

// Driver is a stockdb interface
type Driver interface {
	close() error

	PutMarket(market string) response
	PutOHLC(datum stockdb.OHLC, opt stockdb.Option) response
	PutOHLCs(data []stockdb.OHLC, opt stockdb.Option) response
}
