package clientname

import (
	"github.com/cuigh/auxo/net/rpc"
	"context"
)

type contextKey struct{}

var ctxKey = contextKey{}

func Server() rpc.SFilter {
	return func(next rpc.SHandler) rpc.SHandler {
		return func(c rpc.Context) (r interface{}, err error) {
			c.SetContext(context.WithValue(c.Context(), ctxKey, c.Request().Head.Labels.Get("client.name")))
			r, err = next(c)
			return
		}
	}
}

func FromContext(c context.Context) string {
	return c.Value(ctxKey).(string)
}