package handler

import (
	"encoding/json"
	"knands42/url-shortener/internal/database/repo"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type GetURLResponse struct {
	OriginalUrl string `json:"original_url"`
	ShortUrl    string `json:"short_url"`
}

const (
	URL_TYPE_SHORT    = "short_url"
	URL_TYPE_ORIGINAL = "original_url"
)

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	url := chi.URLParam(r, "url")
	// add default value for url_type query param if not provided
	var urlType string = URL_TYPE_SHORT
	if r.URL.Query().Get("type") != "" {
		urlType = r.URL.Query().Get("type")
	}

	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	var resultDB repo.ShortenedUrl
	var err error
	if urlType == URL_TYPE_SHORT {
		resultDB, err = h.repo.GetByShortUrl(r.Context(), url)
	} else {
		resultDB, err = h.repo.GetByOriginalUrl(r.Context(), url)
	}

	// TODO: Validate error types of sqlc
	if err != nil {
		log.Printf("Failed to get URL: %v", err)
		http.Error(w, "Failed to get URL", http.StatusInternalServerError)
		return
	}

	resp := GetURLResponse{
		OriginalUrl: resultDB.OriginalUrl,
		ShortUrl:    resultDB.ShortUrl,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
