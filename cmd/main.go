package main

import (
	"knands42/url-shortener/internal/database"
	handler "knands42/url-shortener/internal/handlers"
	"knands42/url-shortener/internal/server"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	// TODO: Load env environments

	// Initialize the database
	dbConfig := database.NewDBConfig()
	_, err := dbConfig.Connect()
	if err != nil {
		panic(err)
	}

	// Initialize the handlers
	handlers := handler.NewHandler()

	// Initialize the server
	server.NewServer(r, handlers)
	http.ListenAndServe(":3333", r)
}
