package exception

import (
	"fmt"
	"net/http"
	"reflect"
)

func ReportErr(err error) {
	// 记录日志信息
	fmt.Println(err.Error())
}

func RenderErr(req *http.Request, rspWriter http.ResponseWriter, err error) {
	// 向外部输出错误提示
	msg := ""
	t := reflect.TypeOf(err)
	switch t.Name() {
	case "NotFoundHttpError":
		msg = "404 not found"
	case "MethodNotAllowedHttpError":
		msg = "403 forbiden"
	default:
		msg = "Something went wrong"
	}
	rspWriter.Write([]byte(msg))
}
