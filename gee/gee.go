package gee

import (
	"fmt"
	"net/http"
)

type HandleFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandleFunc
}

func (engine Engine) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else{
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

// 构造方法
func New() *Engine {
	return &Engine{
		router: make(map[string]HandleFunc),
	}
}

// 添加路由
func (engine Engine) addRoute(method string, pattern string, handler HandleFunc) {
	engine.router[method + "-" + pattern] = handler
}

func (engine Engine) GET(pattern string, handler HandleFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine Engine) POST(pattern string, handler HandleFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}