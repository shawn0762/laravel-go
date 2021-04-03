package app

import (
	"app/context"
	"app/exception"
	"app/routes"
	"app/routing"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type app struct {
	Router *routing.Router
}

func NewApp() *app {
	app := &app{Router: routing.NewRouter()}
	// 注册控制器
	app.registerHandlers()
	return app
}

// 封装http server启动
func (app app) Start() {
	log.Fatal(http.ListenAndServe(":9999", app))
}

// 实现Handler接口
func (app app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			e, ok := err.(error)
			if ok {
				exception.ReportErr(e)
				exception.RenderErr(r, w, e)
			}
		}
	}()

	ctx := context.NewContext(w, r)
	route, err := app.Router.GetRoute(ctx.Method, ctx.Uri)
	if err != nil {
		e := errors.New("Something want wrong")
		panic(e)
		return
	}
	// 执行业务逻辑
	route.Handler(ctx)
}

func (app *app) registerHandlers() {
	routes.Group(app.Router)
	fmt.Println(app.Router)
}
