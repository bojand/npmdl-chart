package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", Index)
	router.Get("/chart/*", DrawNPMChart)
	router.Get("/{name}", GetNPMChart)
	router.Get("/{name}/*", GetNPMChart)
	FileServer(router, "/static", "static")

	http.ListenAndServe(":8080", router)
}
