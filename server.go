package main

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/hprose/hprose-golang/rpc"
)

type response struct {
	Success bool        `json:"Success"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

func (response) OnSendHeader(ctx *rpc.HTTPContext) {
	ctx.Response.Header().Set("Access-Control-Allow-Headers", "Authorization")
}

func server() {
	service := rpc.NewHTTPService()
	service.Event = response{}
	service.AddBeforeFilterHandler(func(request []byte, ctx rpc.Context, next rpc.NextFilterHandler) (response []byte, err error) {
		ctx.SetInt64("start", time.Now().UnixNano())
		httpContext := ctx.(*rpc.HTTPContext)
		if httpContext != nil && httpContext.Request.Header.Get("Authorization") == config["http.auth"] {
			ctx.SetBool("authorized", true)
		}
		return next(request, ctx)
	})
	service.AddInvokeHandler(func(name string, args []reflect.Value, ctx rpc.Context, next rpc.NextInvokeHandler) (results []reflect.Value, err error) {
		if openMethods[name] || ctx.GetBool("authorized") {
			results, err = next(name, args, ctx)
		} else {
			resp := response{Message: errHTTPUnauthorized.Error()}
			results = append(results, reflect.ValueOf(resp))
		}
		if logConf.Enable {
			spend := (time.Now().UnixNano() - ctx.GetInt64("start")) / 1000000
			spendInfo := ""
			if spend > 1000 {
				spendInfo = fmt.Sprintf("%vs", spend/1000)
			} else {
				spendInfo = fmt.Sprintf("%vms", spend)
			}
			log(logRequest, fmt.Sprintf("%12s() spend %s", name, spendInfo))
		}
		return
	})
	service.AddMethods(
		[]string{
			"PutOHLC",
			"PutOHLCs",
			"PutOrder",
			"PutOrders",
			"GetStats",
			"GetMarkets",
			"GetSymbols",
			"GetTimeRange",
			"GetPeriodRange",
			"GetOHLCs",
			"GetDepth",
		},
		newInfluxdb(),
		nil,
	)
	http.Handle("/", service)
	http.Handle("/admin/", http.FileServer(http.Dir("")))
	if err := http.ListenAndServe(config["http.bind"], nil); err != nil {
		log(logFatal, "Server error: ", err)
	}
}
