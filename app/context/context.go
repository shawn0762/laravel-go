package context

import "net/http"

type Context struct {
	RspWrite   http.ResponseWriter
	Req        *http.Request
	Method     string
	Uri        string
	RspBuilder *responseBuilder
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		RspWrite:   w,
		Req:        r,
		Method:     r.Method,
		Uri:        r.URL.Path,
		RspBuilder: &responseBuilder{rspWriter: w},
	}
}
