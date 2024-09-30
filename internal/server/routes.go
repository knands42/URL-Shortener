package server

import (
	"github.com/go-chi/chi/v5"
)

func (s *Server) ConfigureRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/shorten", s.handlers.GenerateShortURL)

	return r
}
