package server

import (
	handler "knands42/url-shortener/internal/handlers"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	Router   *chi.Mux
	handlers *handler.Handler
	tracing  trace.Tracer
}

func NewServer(router *chi.Mux, handlers *handler.Handler, tracing trace.Tracer) *Server {
	server := &Server{Router: router, handlers: handlers, tracing: tracing}

	server.DefaultMiddlewares()
	server.DefaultRoutes()
	server.CustomRoutes()

	return server
}
