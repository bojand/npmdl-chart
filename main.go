package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	router := chi.NewRouter()
	router.Get("/", Index)
	router.Get("/chart/*", DrawNPMChart)
	router.Get("/{name}", GetNPMChart)
	FileServer(router, "/static", "static")

	http.ListenAndServe(":8080", router)
}
