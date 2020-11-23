package main

import (
	"fmt"
	"gee"
	"log"
	"net/http"
)

func main() {
	engine := gee.New()

	engine.GET("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "URL.PATH = %q\n", request.URL.Path)
	})

	engine.GET("/hello", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "hello world!\n")
		for k, v := range request.Header {
			fmt.Fprintf(writer, "Header[%q] = %q\n", k, v)
		}
	})

	err := engine.Run(":9999")
	log.Fatal(err)
}
