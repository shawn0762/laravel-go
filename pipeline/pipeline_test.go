package pipeline

import (
	"fmt"
	"testing"
)

func TestPipeline_Then(t *testing.T) {
	p := NewPipeline()

	p1 := func(ps passable, next next) {
		fmt.Println("p1")
		next(ps)
	}

	p2 := func(ps passable, next next) {
		fmt.Println("p2")
		next(ps)
	}

	dis := func() {
		fmt.Println("destination")
	}
	p.Through(p1, p2).Then(dis)
}
