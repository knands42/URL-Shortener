package handler

import (
	"encoding/json"
	"knands42/url-shortener/internal/database/repo"
	"log"
	"net/http"
	"strings"
)

type GetURLResponse struct {
	OriginalUrl string `json:"original_url"`
	ShortUrl    string `json:"short_url"`
}

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	urlQueryParam := r.URL.Query().Get("url")
	if urlQueryParam == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	var err error
	var resultDB repo.ShortenedUrl

	if strings.Contains(urlPath, "/shorten") {
		resultDB, err = h.repo.GetByOriginalUrl(r.Context(), urlQueryParam)
	} else {
		resultDB, err = h.repo.GetByShortUrl(r.Context(), urlQueryParam)
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
