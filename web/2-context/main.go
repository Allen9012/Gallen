/**
  @author: Allen
  @since: 2023/3/22
  @desc: //TODO
**/
package main

import (
	"github/Allen9012/gellen/gellen"
	"net/http"
)

func main() {
	r := gellen.New()
	r.GET("/", func(c *gellen.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	r.GET("/hello", func(c *gellen.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	r.POST("/login", func(c *gellen.Context) {
		c.JSON(http.StatusOK, gellen.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
}
