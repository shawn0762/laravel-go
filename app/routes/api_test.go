package routes

import (
	"app/routing"
	"fmt"
	"testing"
)

func TestGroup(t *testing.T) {
	r := routing.NewRouter()
	Group(r)

	route, ok := r.GetRoute("GET", "/v1/order/b")
	if ok != nil {
		fmt.Println("something wrong")
		return
	}

	fmt.Println(route)
}
