package server

import "github.com/go-chi/chi/v5/middleware"

func (s *Server) DefaultMiddlewares() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	// TODO: Create default error handler
}
