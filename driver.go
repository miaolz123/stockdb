package main

type option struct {
	Market string `json:"market" ini:"market"`
	Symbol string `json:"symbol" ini:"symbol"`
	Period int64  `json:"period" ini:"period"`
}

type ticker struct {
	ID     string  `json:"id"`
	Time   int64   `json:"time"`
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
}

type ohlc struct {
	Time   int64   `json:"time"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

// Driver is a stockdb interface
type Driver interface {
	close() error

	AddMarket(market string) response
	AddOHLC(datum ohlc, opt option) response
	AddOHLCs(data []ohlc, opt option) response
}
