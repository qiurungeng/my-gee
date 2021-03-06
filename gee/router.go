package gee

import (
	"log"
	"net/http"
	"reflect"
	"strings"
)

type router struct {
	roots map[string]*node
	handlers map[string]HandleFunc
}

func newRouter() *router{
	return &router{
		roots: make(map[string]*node),
		handlers: make(map[string]HandleFunc),
	}
}

// 添加路由处理函数
func (r *router) addRoute(method string, pattern string, handler HandleFunc)  {
	parts := patternToParts(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		// 第0层, 根节点为空
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler

	log.Printf("Route %4s - %s , ROUTER:%s", method, pattern, reflect.TypeOf(r))
}

// 从路由数匹配路由节点
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := patternToParts(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := patternToParts(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}


// Context路由处理
func (r router) handle(ctx *Context) {
	node, params := r.getRoute(ctx.Method, ctx.Path)

	// 将原本路由所指向的 handle 函数打包入 context.handlers
	// 这样 handlers 将包含 中间件函数 及 handle 函数
	if node != nil{
		ctx.Params = params
		key := ctx.Method + "-" + node.pattern
		ctx.handlers = append(ctx.handlers, r.handlers[key])
	}else {
		ctx.handlers = append(ctx.handlers, func(ctx *Context) {
			ctx.String(http.StatusNotFound, "404 NOT FOUND:%s\n", ctx.Path)
		})
	}
	// 链式调用
	ctx.Next()
}


// 分隔 Path pattern
func patternToParts(pattern string) []string {
	split := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range split{
		if item != "" {
			parts = append(parts, item)
			//只有最后一项能为 * 通配符
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}