package main

import (
	"github.com/cuigh/auxo/app"
	"github.com/cuigh/auxo/net/rpc"
	"github.com/cuigh/auxo/net/rpc/codec/http"
)

func main() {
	app.Action = func(c *app.Context) {
		app.Run(server())
	}
	app.Start()
}

func server() *rpc.Server {
	//s, _ := rpc.Listen(transport.Address{URL: ":9000"})
	s, err := rpc.AutoServer("test")
	if err != nil {
		panic(err)
	}
	s.Match(http.Matcher, "http")
	//s.Use()
	return s
}

