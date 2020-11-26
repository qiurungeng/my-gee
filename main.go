package main

import (
	"fmt"
	"gee"
	"html/template"
	"log"
	"net/http"
	"time"
)

func main() {
	engine := gee.New()

	// 开启日志，错误恢复
	engine.AddMiddlewares(gee.Logger(), gee.Recovery())

	// 设置模板渲染函数集
	engine.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	// 设置模板文件路径
	engine.LoadHTMLGlob("templates/*")
	// 设置静态文件路径映射
	engine.Static("/assets", "./static")

	// 注册路由
	{
		engine.GET("/", func(ctx *gee.Context) {
			ctx.HTML(http.StatusOK, "css.tmpl", nil)
		})
		engine.GET("/students", func(ctx *gee.Context) {
			ctx.HTML(http.StatusOK, "arr.tmpl", gee.H{
				"title": "gee",
				"stuArr": [2]*student{
					{Name: "Tom", Age: 15},
					{Name: "Jerry", Age: 3},
				},
			})
		})
		engine.GET("/data", func(ctx *gee.Context) {
			ctx.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
				"title": "gee",
				"now": time.Now(),
			})
		})
		engine.GET("/testPanic", func(ctx *gee.Context) {
			//将直接数组越界导致程序崩溃，测试错误恢复
			names := []string{"大大大"}
			ctx.String(http.StatusOK, names[100])
		})
	}

	err := engine.Run(":9999")
	log.Fatal(err)
}

func logMiddlewareForV2() gee.HandleFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		c.Next()
		// 若发生服务内部错误，快速失败，c.index = len(c.handlers)
		if c.StatusCode != http.StatusOK {
			c.Fail(500, "Internal Server Error")
		}
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}


type student struct {
	Name string
	Age int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}
