package main

import (
	"net/http"

	"github.com/hprose/hprose-golang/rpc"
)

type hproseEvent struct{}

func (e hproseEvent) OnSendHeader(ctx *rpc.HTTPContext) {
	ctx.Response.Header().Set("Access-Control-Allow-Headers", "Authorization")
}

func server() {
	service := rpc.NewHTTPService()
	service.Event = hproseEvent{}
	service.AddBeforeFilterHandler(func(request []byte, ctx rpc.Context, next rpc.NextFilterHandler) (response []byte, err error) {
		httpContext := ctx.(*rpc.HTTPContext)
		if httpContext == nil || httpContext.Request.Header.Get("Authorization") != config["http.auth"] {
			httpContext.Response.WriteHeader(http.StatusUnauthorized)
			httpContext.Response.Header().Set("WWW-Authenticate", "Stockdb Server")
			return []byte(errHTTPUnauthorized.Error()), errHTTPUnauthorized
		}
		return next(request, ctx)
	})
	// service.AddInstanceMethods(newInfluxdb())
	if err := http.ListenAndServe(config["http.bind"], service); err != nil {
		log(logFatal, "Server error: ", err)
	}
}
