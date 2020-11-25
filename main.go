package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

func main() {
	engine := gee.New()

	engine.GET("/index", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "<h1>INDEX PAGE</h1>")
	})

	v1 := engine.NewRouterGroup("/v1")
	v1.AddMiddlewares(gee.Logger())
	{
		v1.GET("/", func(ctx *gee.Context) {
			ctx.HTML(http.StatusOK,"<h1>你好啊，这是首页，路径是: /</h1>")
		})
		v1.GET("/hello", func(ctx *gee.Context) {
			ctx.String(http.StatusOK, "你好啊， %s，您访问的地址是：%s。", ctx.Query("name"), ctx.Path)
		})
	}

	v2 := engine.NewRouterGroup("/v2")
	v2.AddMiddlewares(logMiddlewareForV2())
	{
		v2.GET("/hello/:name", func(ctx *gee.Context) {
			ctx.String(http.StatusOK, "你好啊， %s，您访问的地址是：%s。", ctx.Param("name"), ctx.Path)
		})

		v2.GET("/assets/*filepath", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
		})

		v2.POST("/login", func(ctx *gee.Context) {
			ctx.JSON(http.StatusOK, gee.H{
				"username": ctx.PostForm("username"),
				"password": ctx.PostForm("password"),
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
