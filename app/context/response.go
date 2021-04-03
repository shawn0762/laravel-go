package context

import (
	"encoding/json"
	"net/http"
)

type responseBuilder struct {
	rspWriter http.ResponseWriter
}

func (b *responseBuilder) String(statusCode int, content string) {
	b.setContentType("text/plain")
	b.Status(statusCode)
	b.rspWriter.Write([]byte(content))
}

func (b *responseBuilder) Html(statusCode int, content string) {
	b.setContentType("text/html")
	b.Status(statusCode)
	b.rspWriter.Write([]byte(content))
}

func (b *responseBuilder) Json(statusCode int, content interface{}) {
	b.setContentType("application/json")
	b.Status(statusCode)
	encode := json.NewEncoder(b.rspWriter)
	err := encode.Encode(content)
	if err != nil {
		http.Error(b.rspWriter, err.Error(), 500)
		return
	}
}

func (b *responseBuilder) Redirect(to string) {
	b.setContentType("text/html")
	b.SetHeader("Location", to)
	b.Status(302)
}

func (b *responseBuilder) setContentType(value string) {
	b.SetHeader("Content-Type", value)
}

func (b *responseBuilder) SetHeader(key string, value string) {
	b.rspWriter.Header().Set(key, value)
}

// 这个方法的调用时机必须在:
// 1.ResponseWriter开始写入内容之前调用
// 2.请求头设置完成之后
func (b *responseBuilder) Status(statusCode int) {
	b.rspWriter.WriteHeader(statusCode)
}
