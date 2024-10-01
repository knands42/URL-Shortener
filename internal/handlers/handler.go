package handler

import "knands42/url-shortener/internal/database/repo"

type Handler struct {
	repo *repo.Queries
}

func NewHandler(repo *repo.Queries) *Handler {
	return &Handler{
		repo: repo,
	}
}
