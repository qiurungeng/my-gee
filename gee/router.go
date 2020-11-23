package gee

import (
	"log"
	"net/http"
)

type router struct {
	handlers map[string]HandleFunc
}

func newRouter() *router{
	return &router{
		handlers: make(map[string]HandleFunc),
	}
}

// 添加路由处理函数
func (r router) addRoute(method string, pattern string, handler HandleFunc)  {
	log.Printf("Route %4s - %s", method, pattern)
	r.handlers[method + "-" + pattern] = handler
}

// 路由处理context
func (r router) handle(ctx *Context) {
	key := ctx.Method + "-" + ctx.Path
	if handler, ok := r.handlers[key]; ok{
		handler(ctx)
	} else{
		ctx.String(http.StatusNotFound, "404 NOT FOUND:%s\n", ctx.Path)
	}
}

