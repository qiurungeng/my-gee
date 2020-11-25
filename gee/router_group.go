package gee

import (
	"net/http"
	"path"
)

//分组控制(Group Control)是 Web 框架应提供的基础功能之一
//所谓分组，是指路由的分组。如果没有路由分组，我们需要针对每一个路由进行控制
//但是真实的业务场景中，往往某一组路由需要相似的处理
//RouterGroup 路由分组
type RouterGroup struct {
	//该分组路由前缀
	prefix string
	//中间件
	middlewares []HandleFunc
	//父分组
	parent *RouterGroup
	//所属 Engine 实例
	engine *Engine
}

// NewRouterGroup: 为当前路由分组创建一个子路由分组
// 所有路由分组共享一个 Engine 实例
func (group *RouterGroup) NewRouterGroup(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.allGroups = append(engine.allGroups, newGroup)
	return newGroup
}

// 为 RouterGroup 添加一个或多个中间件
func (group *RouterGroup) AddMiddlewares(middlewares ...HandleFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// 为 RouterGroup 添加路由
func (group *RouterGroup) addRoute(method string, subPattern string, handler HandleFunc) {
	pattern := group.prefix + subPattern
	//log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandleFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandleFunc) {
	group.addRoute("POST", pattern, handler)
}



/******************************静态资源处理**********************************/

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandleFunc {
	absolutePath := path.Join(group.prefix, relativePath)

	// fileServer: HandleFunc
	// 去除请求URL的前缀: absolutePath
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// 用户可以将磁盘上的某个文件夹root映射到路由relativePath, 例如:
// r := gee.New()
// r.Static("/assets", "/usr/blog/static")
// 用户访问localhost:9999/assets/js/my.js，
// 最终返回/usr/blog/static/js/my.js。
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}