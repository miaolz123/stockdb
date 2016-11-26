package main

type response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type driver interface {
	reconnect()
	close() error
	check() error

	AddMarket(market string) response
	AddOHLC(market, symbol string, datum ohlc) response
	AddOHLCs(market, symbol string, data []ohlc) response
}
