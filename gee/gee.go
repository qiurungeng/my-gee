package gee

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

type HandleFunc func(ctx *Context)

type Engine struct {
	*RouterGroup						//Engine 的根路由分组，默认自带的
	router 			*router
	allGroups 		[]*RouterGroup 		//Engine 上的所有路由分组
	htmlTemplates 	*template.Template	//for html render
	funcMap 		template.FuncMap	//for html render
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	// 为context添加其所属路由分组所适用的所有中间件
	var middlewares []HandleFunc
	for _, group := range engine.allGroups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}

// 构造方法
func New() *Engine {
	engine := &Engine{ router: newRouter() }
	engine.RouterGroup = &RouterGroup{engine: engine} //Engine 的根路由分组，默认自带的
	engine.allGroups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 启动服务
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}


/**************** Template 渲染 相关方法 *******************
 */

// 自定义模板渲染函数
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.
			New("").
			Funcs(engine.funcMap).	// 函数集，可在模板渲染过程使用
			ParseGlob(pattern))		// pattern 匹配文件
}

//----------- Util Func -----------
func Logger() HandleFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}