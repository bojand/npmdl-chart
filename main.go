package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	router := fasthttprouter.New()
	router.GET("/chart/*name", DrawNPMChart)

	fasthttp.ListenAndServe(":8080", router.Handler)
}
