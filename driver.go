package main

import (
	"github.com/miaolz123/stockdb/stockdb"
)

// Driver is a stockdb interface
type Driver interface {
	close() error

	PutOHLC(datum stockdb.OHLC, opt stockdb.Option) response
	PutOHLCs(data []stockdb.OHLC, opt stockdb.Option) response
	PutOrder(datum stockdb.Order, opt stockdb.Option) response
	PutOrders(data []stockdb.Order, opt stockdb.Option) response
	GetStats() response
	GetMarkets() response
	GetSymbols(market string) response
	GetTimeRange(opt stockdb.Option) response
	GetPeriodRange(opt stockdb.Option) response
	GetOHLCs(opt stockdb.Option) response
	GetDepth(opt stockdb.Option) response
}
