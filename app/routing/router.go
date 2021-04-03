package routing

import (
	"app/exception"
	"errors"
	"strings"
)

/**
 * 功能要求：
 * 1、路由可选参数：/a/{id}
 * 2、路由必选参数：/a/{id?}
 * 3、命名路由
 */

// 路由注册、路由匹配
type Router struct {
	// 存放所有路由信息：按照method分组
	routes map[string]*node

	groupIndex int
	groupStack map[int]*group
}

func NewRouter() *Router {
	return &Router{
		routes:     make(map[string]*node),
		groupIndex: 0,
		groupStack: make(map[int]*group),
		//actions: nil,
	}
}

func (r *Router) Get(uri string, handler interface{}) {
	r.registerRoute("GET", uri, handler)
}

func (r *Router) Post(uri string, handler interface{}) {
	r.registerRoute("POST", uri, handler)
}

func (r *Router) Put(uri string, handler interface{}) {
	r.registerRoute("PUT", uri, handler)
}

// 根据层层嵌套的group参数，初始化一个完整的route对象
func (r *Router) registerRoute(method string, uri string, handler interface{}) {
	// 读取上一层的group信息
	prevGroup, ok := r.groupStack[r.groupIndex]
	if !ok {
		prevGroup = &group{}
	}

	currGroup := &group{}
	hdl, ok := handler.(Handler)
	if !ok {

		attrs, ok := handler.(map[string]interface{})
		if !ok {
			// @todo 有更优雅的方式?
			panic(errors.New("Invalid action"))
		}

		currGroup, hdl = r.parseHandler(attrs)
	}

	// 拼上最后一段
	pattern := prevGroup.prefix + "/" + strings.Trim(uri, "/")
	as := prevGroup.as + "." + currGroup.as

	// 这里要把路由注册到Router
	route := &Route{
		pattern: strings.Trim(pattern, "/"),
		method:  method,
		as:      strings.Trim(as, "."),
		Handler: hdl,
	}
	r.addRoute(route)
}

func (r *Router) parseHandler(attrs map[string]interface{}) (*group, Handler) {
	tmpAs := attrs["as"]
	as, ok := tmpAs.(string)
	if !ok {
		as = ""
	}
	g := &group{as: as}

	uses := attrs["uses"]
	h, ok := uses.(Handler)
	if !ok {
		// @todo 有更优雅的方式?
		panic(errors.New("Action not found"))
	}
	return g, h
}

// 将route对象插入trie树中
func (r *Router) addRoute(route *Route) {
	method := route.method
	root, ok := r.routes[method]
	if !ok {
		root = NewRoot()
		r.routes[method] = root
	}

	path := route.pattern
	parts := parsePath(path)
	root.insert(route, path, parts, 0)
}

// 根据请求uri，匹配一个route
func (r Router) GetRoute(method string, path string) (*Route, error) {
	root, ok := r.routes[method]
	if !ok {
		// @todo 这里可以有更优雅的错误响应？
		panic(exception.MethodNotAllowedHttpError{})
	}
	//path = strings.Trim(path, "/")
	parts := parsePath(path)
	node := root.Search(parts, 0)
	if node == nil {
		// @todo 这里可以有更优雅的错误响应？
		panic(exception.NotFoundHttpError{})
	}

	// @todo 获取路由参数
	// 		 e.g. /a/b/1 匹配到/a/{type}/{id}
	// 		 则有type=b    id=1

	return node.Route, nil
}

func parsePath(path string) []string {
	path = strings.Trim(path, "/")
	return strings.Split(path, "/")
}

// 定义Group方法第一个参数的格式
type group struct {
	prefix string
	as     string
}

func (r *Router) Group(attrs map[string]interface{}, handler GroupHandler) {

	r.updateGroupStack(attrs)

	handler(r)

	delete(r.groupStack, r.groupIndex)
	r.groupIndex -= 1
}

func (r *Router) updateGroupStack(attr map[string]interface{}) {
	// 获取每一层group的参数
	// 然后将参数拼接父group的参数
	// 将拼接好的参数较费下一层调用
	tmp, exists := attr["as"]
	as, isString := tmp.(string)
	if !exists || !isString {
		as = ""
	}

	tmp, exists = attr["prefix"]
	prefix, isString := tmp.(string)
	if !exists || !isString {
		prefix = ""
	}
	prefix = strings.Trim(prefix, "/")

	// 每次执行group方法,都会在groupStack追加一个新group
	// 供group方法里面的方法使用
	// 当前group方法调用完毕则移除这个新group
	// 这样兄弟group方法的调用就能取到父group的对象,继续合并需要的参数
	prevStack, ok := r.groupStack[r.groupIndex]
	if !ok {
		prevStack = &group{}
	}

	newStack := &group{
		prefix: strings.Trim(prevStack.prefix+"/"+prefix, "/"),
		as:     strings.Trim(prevStack.as+"."+as, "."),
	}
	r.groupIndex += 1
	r.groupStack[r.groupIndex] = newStack
}

func mergeToGroup(g group, attr map[string]interface{}) {

}

// 这里定义一个方法类型，限定group方法的第二个参数
type GroupHandler func(route *Router)
