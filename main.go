package main

import (
	"app"
)

func main() {
	// 启动服务容器
	container := app.NewApp()
	container.Start()
}
