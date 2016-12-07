package stockdb

import (
	"encoding/base64"
	"net/http"

	"github.com/hprose/hprose-golang/io"
	"github.com/hprose/hprose-golang/rpc"
)

// Option is a request option
type Option struct {
	Market        string `json:"Market" ini:"Market"`
	Symbol        string `json:"Symbol" ini:"Symbol"`
	Period        int64  `json:"Period" ini:"Period"`
	BeginTime     int64  `json:"BeginTime" ini:"BeginTime"`
	EndTime       int64  `json:"EndTime" ini:"EndTime"`
	InvalidPolicy string `json:"InvalidPolicy" ini:"InvalidPolicy"`
}

// BaseResponse is base response struct
type BaseResponse struct {
	Success bool        `json:"Success"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

// Ticker is an order record struct
type Ticker struct {
	ID     string  `json:"ID"`
	Time   int64   `json:"Time"`
	Price  float64 `json:"Price"`
	Amount float64 `json:"Amount"`
	Type   string  `json:"Type"`
}

// OHLC is a candlestick struct
type OHLC struct {
	Time   int64   `json:"Time"`
	Open   float64 `json:"Open"`
	High   float64 `json:"High"`
	Low    float64 `json:"Low"`
	Close  float64 `json:"Close"`
	Volume float64 `json:"Volume"`
}

// OrderBook struct
type OrderBook struct {
	Price  float64 `json:"Price"`
	Amount float64 `json:"Amount"`
}

// Depth struct
type Depth struct {
	Bids []OrderBook `json:"Bids"`
	Asks []OrderBook `json:"Asks"`
}

// TimeRangeResponse is TimeRange response struct
type TimeRangeResponse struct {
	Success bool     `json:"Success"`
	Message string   `json:"Message"`
	Data    [2]int64 `json:"Data"`
}

// OHLCResponse is OHLC response struct
type OHLCResponse struct {
	Success bool   `json:"Success"`
	Message string `json:"Message"`
	Data    []OHLC `json:"Data"`
}

// DepthResponse is market depth response struct
type DepthResponse struct {
	Success bool   `json:"Success"`
	Message string `json:"Message"`
	Data    Depth  `json:"Data"`
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
	GetDepth     func(opt Option) DepthResponse
}

// New can create a StockDB Client
func New(uri, auth string) (client *Client) {
	io.Register(Option{}, "Option", "json")
	io.Register(BaseResponse{}, "BaseResponse", "json")
	io.Register(Ticker{}, "Ticker", "json")
	io.Register(OHLC{}, "OHLC", "json")
	io.Register(OHLCResponse{}, "OHLCResponse", "json")
	io.Register(OrderBook{}, "OrderBook", "json")
	io.Register(Depth{}, "Depth", "json")
	io.Register(DepthResponse{}, "DepthResponse", "json")
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
