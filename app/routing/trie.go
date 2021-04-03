package routing

import (
	"regexp"
	"strings"
)

type node struct {
	pattern string // 只有在叶子节点才会有值

	part string

	children []*node

	Route *Route // 只有在叶子节点才会有值

	isWide bool

	isOption bool // 是否是可选的
}

func Match(pattern string, root *node) *node {
	parts := strings.Split(pattern, "/")
	return root.Search(parts, 0)
}

func NewRoot() *node {
	return &node{
		pattern:  "/",
		part:     "",
		children: nil,
		Route:    nil,
		isWide:   false,
	}
}

func (n *node) Search(parts []string, height int) *node {
	// 因为当前节点是有上一层match children出来的，所以肯定能匹配路由
	// 终止递归的条件就是：是否已经匹配到路由的最后一段
	// 一旦匹配到最后一段，此时肯定要有一个结果：要么是匹配成功，要么是匹配失败

	// 只要深度到达最后一段，就可以开始有结果
	//
	l := len(parts)
	if l <= height && n.pattern != "" {
		return n
	}

	var part string
	if height < l {
		part = parts[height]
	}

	children := n.matchChildren(part)
	for _, child := range children {
		ret := child.Search(parts, height+1)
		if ret != nil {
			return ret
		}
	}
	// 没有一个能够匹配上，就不用继续找了
	return nil
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}
	return nil
}

// 查找所有能够匹配路由的children
func (n *node) matchChildren(part string) []*node {
	children := make([]*node, 0)
	for _, child := range n.children {
		if part != "" {
			if part == child.part || child.isWide {
				children = append(children, child)
			}
		} else {
			// 如果part为空，说明当前深度已经比实际uri大
			// 此时后面的节点必须是可选参数
			// 比如：/a 可以匹配 /a/{id?}/{type?}
			if child.isOption {
				children = append(children, child)
			}
		}
	}
	return children
}

//
func (n *node) insert(route *Route, pattern string, parts []string, height int) {
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		// 说明当前层没有这个节点，插入一个新的节点
		isWide, ok := regexp.Match(`^\{.+\}$`, []byte(part))
		if ok != nil {
			isWide = false
		}

		isOption, ok := regexp.Match(`^\{.+\?\}$`, []byte(part))
		if ok != nil {
			isOption = false
		}

		child = &node{
			pattern:  "",
			part:     part,
			children: nil,
			Route:    nil,
			isWide:   isWide,
			isOption: isOption,
		}
		n.children = append(n.children, child)

		if len(parts)-1 == height {
			child.pattern = pattern
			child.Route = route
			return
		}
	}
	child.insert(route, pattern, parts, height+1)
}
