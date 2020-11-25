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

	engine.AddMiddlewares(gee.Logger())

	engine.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	engine.LoadHTMLGlob("templates/*")
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
