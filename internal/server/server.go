package server

import (
	handler "knands42/url-shortener/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router   *chi.Mux
	handlers *handler.Handler
}

func NewServer(router *chi.Mux) *Server {
	server := &Server{router: router}

	server.DefaultMiddlewares()
	server.DefaultRoutes()
	server.ApplicationRoutes()

	// TOOD: Add cors
	// TODO: Add default headers
	return server
}

func (s *Server) DefaultMiddlewares() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	// TODO: Create default error handler
}

func (s *Server) DefaultRoutes() {
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
}

func (s *Server) ApplicationRoutes() http.Handler {
	s.router.Mount("/api/v1", s.ConfigureRoutes())
	return s.router
}
