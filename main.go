package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	router := fasthttprouter.New()
	router.GET("/", Index)
	router.GET("/chart/*name", DrawNPMChart)
	// router.GET("/:name", GetNPMChart)
	router.ServeFiles("/static/*filepath", "static")

	fasthttp.ListenAndServe(":8080", router.Handler)
}
