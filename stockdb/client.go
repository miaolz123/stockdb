package stockdb

import (
	"encoding/base64"
	"net/http"

	"github.com/hprose/hprose-golang/rpc"
)

// Option is a request option
type Option struct {
	Market string `json:"market" ini:"market"`
	Symbol string `json:"symbol" ini:"symbol"`
	Period int64  `json:"period" ini:"period"`
}

// BaseResponse is base response struct
type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Ticker is an order record struct
type Ticker struct {
	ID     string  `json:"id"`
	Time   int64   `json:"time"`
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
}

// OHLC is a candlestick struct
type OHLC struct {
	Time   int64   `json:"time"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

// Client of StockDB
type Client struct {
	uri    string
	auth   string
	Hprose *rpc.HTTPClient

	PutMarket func(market string) BaseResponse
	PutOHLC   func(datum OHLC, opt Option) BaseResponse
	PutOHLCs  func(data []OHLC, opt Option) BaseResponse
}

// New can create a StockDB Client
func New(uri, auth string) (client *Client) {
	client = &Client{
		uri:    uri,
		Hprose: rpc.NewHTTPClient(uri),
	}
	if auth != "" {
		client.Hprose.Header = make(http.Header)
		client.Hprose.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
	}
	client.Hprose.UseService(&client)
	return
}
