package handler

import (
	"knands42/url-shortener/internal/database/repo"

	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	repo    *repo.Queries
	tracing trace.Tracer
}

const (
	URL_TYPE_SHORT    = "short_url"
	URL_TYPE_ORIGINAL = "original_url"
)

func NewHandler(repo *repo.Queries, tracing trace.Tracer) *Handler {
	return &Handler{
		repo:    repo,
		tracing: tracing,
	}
}

func (h *Handler) extractHashFromUrl(url string) string {
	return url[len(url)-7:]
}
