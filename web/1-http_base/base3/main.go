/**
  @author: Allen
  @since: 2023/3/22
  @desc: //TODO
**/
package main

import (
	"fmt"
	"github/Allen9012/example/gellen"
	"net/http"
)

func main() {
	r := gellen.New()
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	r.Run(":9090")
}
