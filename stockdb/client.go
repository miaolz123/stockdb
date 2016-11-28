package stockdb

import (
	"encoding/base64"
	"net/http"

	"github.com/hprose/hprose-golang/io"
	"github.com/hprose/hprose-golang/rpc"
)

// Option is a request option
type Option struct {
	Market        string `json:"market" ini:"market"`
	Symbol        string `json:"symbol" ini:"symbol"`
	Period        int64  `json:"period" ini:"period"`
	BeginTime     int64  `json:"beginTime" ini:"beginTime"`
	EndTime       int64  `json:"endTime" ini:"endTime"`
	InvalidPolicy string `json:"invalidPolicy" ini:"invalidPolicy"`
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

// TimeRangeResponse is TimeRange response struct
type TimeRangeResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    [2]int64 `json:"data"`
}

// OHLCResponse is OHLC response struct
type OHLCResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []OHLC `json:"data"`
}

// Client of StockDB
type Client struct {
	uri    string
	auth   string
	Hprose *rpc.HTTPClient

	PutOHLC      func(datum OHLC, opt Option) BaseResponse
	PutOHLCs     func(data []OHLC, opt Option) BaseResponse
	GetTimeRange func(opt Option) TimeRangeResponse
	GetOHLCs     func(opt Option) OHLCResponse
}

// New can create a StockDB Client
func New(uri, auth string) (client *Client) {
	io.Register(Option{}, "Option", "json")
	io.Register(BaseResponse{}, "BaseResponse", "json")
	io.Register(Ticker{}, "Ticker", "json")
	io.Register(OHLC{}, "OHLC", "json")
	io.Register(OHLCResponse{}, "OHLCResponse", "json")
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
