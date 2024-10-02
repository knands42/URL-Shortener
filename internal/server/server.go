package server

import (
	handler "knands42/url-shortener/internal/handlers"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Router   *chi.Mux
	handlers *handler.Handler
}

func NewServer(router *chi.Mux, handlers *handler.Handler) *Server {
	server := &Server{Router: router, handlers: handlers}

	server.DefaultMiddlewares()
	server.DefaultRoutes()
	server.CustomRoutes()

	return server
}
