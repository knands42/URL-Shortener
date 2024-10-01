package server

import (
	handler "knands42/url-shortener/internal/handlers"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	router   *chi.Mux
	handlers *handler.Handler
}

func NewServer(router *chi.Mux, handlers *handler.Handler) *Server {
	server := &Server{router: router, handlers: handlers}

	server.DefaultMiddlewares()
	server.DefaultRoutes()
	server.CustomRoutes()

	// TODO: Add cors
	// TODO: Add default headers
	return server
}
