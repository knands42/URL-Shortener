package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) DefaultRoutes() {
	s.Router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	s.Router.Get("/swagger/*", httpSwagger.WrapHandler)
}

func (s *Server) CustomRoutes() http.Handler {
	s.Router.Mount("/api/v1", s.prefixedAPIs())
	return s.Router
}

func (s *Server) prefixedAPIs() chi.Router {
	r := chi.NewRouter()

	r.Post("/shorten", s.handlers.GenerateShortURL)
	r.Get("/url/{url}", s.handlers.GetOriginalUrl)
	r.Get("/url/{url}/metadata", s.handlers.GetMetadata)
	r.Delete("/url", s.handlers.DeleteURL)

	return r
}
