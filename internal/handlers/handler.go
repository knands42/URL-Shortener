package handler

import (
	"encoding/json"
	"knands42/url-shortener/internal/database/repo"
	"knands42/url-shortener/internal/utils"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	repo    *repo.Queries
	cache   *redis.Client
	tracing trace.Tracer
}

const (
	URL_TYPE_SHORT    = "short_url"
	URL_TYPE_ORIGINAL = "original_url"
)

func NewHandler(repo *repo.Queries, cache *redis.Client, tracing trace.Tracer) *Handler {
	return &Handler{
		repo:    repo,
		cache:   cache,
		tracing: tracing,
	}
}

func (h *Handler) extractHashFromUrl(url string) string {
	return url[len(url)-7:]
}

func notFound(w http.ResponseWriter, err error, msg string) {
	errorResponse := utils.NotFoundErrorResponse{
		Status:  http.StatusNotFound,
		Message: msg,
	}
	log.Printf(msg + " - " + err.Error())
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(errorResponse)
}
