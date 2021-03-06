package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req *http.Request
	// request info
	Path string
	Method string
	Params map[string]string	// 提供对路由参数的访问
	// response info
	StatusCode 	int
	// handlers = middlewares + 原本路由所指向的函数
	handlers []HandleFunc
	index    int
	// Engine
	engine *Engine
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

// Constructor
func newContext(w http.ResponseWriter, req *http.Request) *Context{
	return &Context{
		Writer: w,
		Req: req,
		Path: req.URL.Path,
		Method: req.Method,
		index: -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s ; c.index++{
		c.handlers[c.index](c)
	}
}

// 获取POST请求字段
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 获取GET请求字段
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 设置 response status code 响应状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 设置 response header 响应头kv字段
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 设置 response body 具体内容(字符串)
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain;charset=utf-8")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 设置 response body 具体内容(json)
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json;charset=utf-8")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// 设置 response body 具体内容([]byte])
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// 设置 response body 具体内容(html字符串)
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html;charset=utf-8")	// 注意！utf-8
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}

// Fail: 失败 Response 统一处理
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}
