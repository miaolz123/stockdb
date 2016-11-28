package main

import (
	"net/http"
	"reflect"

	"github.com/hprose/hprose-golang/rpc"
)

type response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (response) OnSendHeader(ctx *rpc.HTTPContext) {
	ctx.Response.Header().Set("Access-Control-Allow-Headers", "Authorization")
}

func server() {
	service := rpc.NewHTTPService()
	service.Event = response{}
	service.AddBeforeFilterHandler(func(request []byte, ctx rpc.Context, next rpc.NextFilterHandler) (response []byte, err error) {
		httpContext := ctx.(*rpc.HTTPContext)
		if httpContext != nil && httpContext.Request.Header.Get("Authorization") == config["http.auth"] {
			ctx.SetBool("authorized", true)
		}
		return next(request, ctx)
	})
	service.AddInvokeHandler(func(name string, args []reflect.Value, ctx rpc.Context, next rpc.NextInvokeHandler) (results []reflect.Value, err error) {
		if ctx.GetBool("authorized") {
			return next(name, args, ctx)
		}
		resp := response{Message: errHTTPUnauthorized.Error()}
		results = append(results, reflect.ValueOf(resp))
		return
	})
	service.AddMethods(
		[]string{
			"PutOHLC",
			"PutOHLCs",
			"GetTimeRange",
			"GetOHLCs",
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
