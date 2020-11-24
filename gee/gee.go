package gee

import (
	"net/http"
)

type HandleFunc func(ctx *Context)

type Engine struct {
	//Engine 的根路由分组，默认自带的
	*RouterGroup
	router *router
	//Engine 上的所有路由分组
	allGroups []*RouterGroup
}

func (engine Engine) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	c := newContext(w, req)
	engine.router.handle(c)
}

// 构造方法
func New() *Engine {
	engine := &Engine{ router: newRouter() }
	engine.RouterGroup = &RouterGroup{engine: engine} //Engine 的根路由分组，默认自带的
	engine.allGroups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 添加路由
func (engine Engine) addRoute(method string, pattern string, handler HandleFunc) {
	//fmt.Println(reflect.TypeOf(engine))
	//engine.routerGroupPrivateMethod("add route: " + method + pattern)
	engine.router.addRoute(method, pattern, handler)
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