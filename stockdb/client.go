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

// OHLC is a candlestick struct
type OHLC struct {
	Time   int64   `json:"Time"`
	Open   float64 `json:"Open"`
	High   float64 `json:"High"`
	Low    float64 `json:"Low"`
	Close  float64 `json:"Close"`
	Volume float64 `json:"Volume"`
}

// Order is an order record struct
type Order struct {
	ID     string  `json:"ID"`
	Time   int64   `json:"Time"`
	Price  float64 `json:"Price"`
	Amount float64 `json:"Amount"`
	Type   string  `json:"Type"`
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

// BaseResponse is base response struct
type BaseResponse struct {
	Success bool        `json:"Success"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

// Stats is stats struct
type Stats struct {
	Market string `json:"Market"`
	Record int64  `json:"Record"`
	Disk   int64  `json:"Disk"`
}

// StatsResponse is Stats response struct
type StatsResponse struct {
	Success bool    `json:"Success"`
	Message string  `json:"Message"`
	Data    []Stats `json:"Data"`
}

// StringsResponse is Strings response struct
type StringsResponse struct {
	Success bool     `json:"Success"`
	Message string   `json:"Message"`
	Data    []string `json:"Data"`
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
	PutOrder     func(datum Order, opt Option) BaseResponse
	PutOrders    func(data []Order, opt Option) BaseResponse
	GetStats     func() StatsResponse
	GetMarkets   func() StringsResponse
	GetSymbols   func(market string) StringsResponse
	GetTimeRange func(opt Option) TimeRangeResponse
	GetOHLCs     func(opt Option) OHLCResponse
	GetDepth     func(opt Option) DepthResponse
}

func init() {
	io.Register(Option{}, "Option", "json")
	io.Register(OHLC{}, "OHLC", "json")
	io.Register(Order{}, "Order", "json")
	io.Register(OrderBook{}, "OrderBook", "json")
	io.Register(Depth{}, "Depth", "json")
	io.Register(BaseResponse{}, "BaseResponse", "json")
	io.Register(Stats{}, "Stats", "json")
	io.Register(StatsResponse{}, "StatsResponse", "json")
	io.Register(StringsResponse{}, "StringsResponse", "json")
	io.Register(TimeRangeResponse{}, "TimeRangeResponse", "json")
	io.Register(OHLCResponse{}, "OHLCResponse", "json")
	io.Register(DepthResponse{}, "DepthResponse", "json")
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
