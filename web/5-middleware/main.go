/**
  @author: Allen
  @since: 2023/3/22
  @desc: //TODO
**/
package main

/*
(1) global 5-middleware Logger
$ curl http://localhost:9999/
<h1>Hello Gee</h1>

>>> log
2019/08/17 01:37:38 [200] / in 3.14µs
*/

/*
(2) global + 4-group 5-middleware
$ curl http://localhost:9999/v2/hello/geektutu
{"message":"Internal Server Error"}

>>> log
2019/08/17 01:38:48 [200] /v2/hello/geektutu in 61.467µs for 4-group v2
2019/08/17 01:38:48 [200] /v2/hello/geektutu in 281µs
*/

import (
	"github/Allen9012/gellen/gellen"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gellen.HandlerFunc {
	return func(c *gellen.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for 4-group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gellen.New()
	r.Use(gellen.Logger()) // global midlleware
	r.GET("/", func(c *gellen.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 4-group 5-middleware
	{
		v2.GET("/hello/:name", func(c *gellen.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	r.Run(":9999")
}
