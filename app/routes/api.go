package routes

import (
	"app/context"
	"app/routing"
)

func Group(r *routing.Router) {
	groupAttr := map[string]interface{}{
		"prefix": "/v1",
		"as":     "pandora.admin",
	}
	r.Group(groupAttr, func(r *routing.Router) {
		r.Get("/user", routing.Handler(user))

		groupAttr := map[string]interface{}{"prefix": "/order", "as": "order."}
		r.Group(groupAttr, func(r *routing.Router) {
			r.Get("info", routing.Handler(order))
			r.Get("/info/{id}", routing.Handler(order2))
			r.Get("/b/{id?}", routing.Handler(order3))
		})
	})
}

func user(ctx *context.Context) {
	ctx.RspBuilder.Html(200, "/v1/user")
}
func order(ctx *context.Context) {
	ctx.RspBuilder.Html(200, "/v1/order/info")
}
func order2(ctx *context.Context) {
	ctx.RspBuilder.Html(200, "/v1/order/info/{id}")
}
func order3(ctx *context.Context) {
	ctx.RspBuilder.Html(200, "/v1/order/b/{id?}")
}
