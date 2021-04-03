package routing

import (
	"app/context"
)

type Route struct {
	pattern string

	method string

	// 路由命名
	as string

	// 业务逻辑处理器
	Handler Handler
}

type Handler func(ctx *context.Context)
