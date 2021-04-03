package pipeline

type pipeline struct {
	pass passable

	pipes []pipe

	destination destination
}

type passable func()
type next func(passable)

type executable func(passable)

type pipe func(passable, next)

type destination func()

func (pipeline *pipeline) Send(pass passable) *pipeline {
	pipeline.pass = pass
	return pipeline
}

func (pipeline *pipeline) Through(pipes ...pipe) *pipeline {
	pipeline.pipes = pipes
	return pipeline
}

func (pipeline *pipeline) Then(dis destination) {
	pipeline.destination = dis
	pipeline.execute()
}

func (pipeline *pipeline) execute() {
	callable := func(pass passable) {
		pipeline.destination()
	}

	// 洋葱模型：
	// 将所有管道一层一层包起来
	// 最先执行的在最外层，最后执行的在最里层
	// 中心则是最终要执行的逻辑
	l := len(pipeline.pipes)
	for i := l - 1; i >= 0; i-- {
		p := pipeline.pipes[i]
		callable = func(carry next, p pipe) executable {
			c := func(pass passable) {
				p(pass, carry)
			}
			return c
		}(callable, p)
	}
	callable(pipeline.pass)
}

func NewPipeline() *pipeline {
	return &pipeline{}
}
