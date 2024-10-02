package server

import (
	"encoding/json"
	"knands42/url-shortener/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) DefaultMiddlewares() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(defaultErrorHandler)
}

func defaultErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(utils.ErrorResponse{
					Status:  http.StatusInternalServerError,
					Message: "Internal Server Error",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
