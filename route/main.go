/**
  @author: Allen
  @since: 2023/3/23
  @desc: //TODO
**/
package route

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
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *gellen.Context) {
		// expect /hello/geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *gellen.Context) {
		c.JSON(http.StatusOK, gellen.H{"filepath": c.Param("filepath")})
	})

	r.Run(":9999")
}
