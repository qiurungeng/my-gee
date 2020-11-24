package gee

import (
	"fmt"
	"reflect"
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

// 为该路由分组添加路由
func (group *RouterGroup) addRoute(method string, subPattern string, handler HandleFunc) {
	pattern := group.prefix + subPattern
	//log.Printf("Route %4s - %s", method, pattern)
	group.engine.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandleFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandleFunc) {
	group.addRoute("POST", pattern, handler)
}

func (group RouterGroup) routerGroupPrivateMethod(str string) {
	fmt.Printf("an instance of %s calls routerGroupPrivateMethod%s\n", reflect.TypeOf(group), str)
}