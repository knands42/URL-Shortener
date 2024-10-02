package handler

import "knands42/url-shortener/internal/database/repo"

type Handler struct {
	repo *repo.Queries
}

const (
	URL_TYPE_SHORT    = "short_url"
	URL_TYPE_ORIGINAL = "original_url"
)

func NewHandler(repo *repo.Queries) *Handler {
	return &Handler{
		repo: repo,
	}
}
