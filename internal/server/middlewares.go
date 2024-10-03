package server

import (
	"encoding/json"
	"knands42/url-shortener/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.opentelemetry.io/otel/trace"
)

func (s *Server) DefaultMiddlewares() {
	s.Router.Use(tracingMiddleware(s.tracing))
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(defaultErrorHandler)
	s.Router.Use(corsMiddleware())
}

// Capture any unhandled errors and return a 500 status code
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

// This will allow every origin to make requests to the server, including the swagger UI
func corsMiddleware() func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}

func tracingMiddleware(tracing trace.Tracer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx, span := tracing.Start(r.Context(), "TracingMiddleware")
			defer span.End()

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
