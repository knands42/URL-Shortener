package main

import (
	"knands42/url-shortener/internal/server"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// Initialize the server
	server.NewServer(r)
	http.ListenAndServe(":3333", r)
}
