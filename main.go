package main

import (
	"github.com/cuigh/auxo/app"
	"github.com/cuigh/auxo/net/rpc"
	"github.com/cuigh/auxo/net/rpc/codec/http"
	"github.com/dongdongking008/softtrans/service"
	"github.com/dongdongking008/softtrans/util/clientname"
)

func main() {
	app.Action = func(c *app.Context) {
		app.Run(server())
	}
	app.Start()
}

func server() *rpc.Server {
	//s, _ := rpc.Listen(transport.Address{URL: ":9000"})
	s, err := rpc.AutoServer("softtrans.coordinator")
	if err != nil {
		panic(err)
	}
	s.Match(http.Matcher, "http")
	//s.Use()
	s.Use(clientname.Server())
	s.RegisterService("TCCService", service.TCCService{})
	return s
}

