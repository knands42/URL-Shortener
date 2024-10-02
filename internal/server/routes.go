package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) DefaultRoutes() {
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	s.router.Get("/swagger/*", httpSwagger.WrapHandler)
}

func (s *Server) CustomRoutes() http.Handler {
	s.router.Mount("/api/v1", s.prefixedAPIs())
	return s.router
}

func (s *Server) prefixedAPIs() chi.Router {
	r := chi.NewRouter()

	r.Post("/shorten", s.handlers.GenerateShortURL)
	r.Get("/shorten", s.handlers.GetURL)
	r.Delete("/shorten", s.handlers.DeleteURL)

	r.Get("/url", s.handlers.GetURL)
	r.Delete("/url", s.handlers.DeleteURL)

	return r
}
