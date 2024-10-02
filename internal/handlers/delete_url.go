package handler

import (
	"log"
	"net/http"
	"strings"
)

func (h *Handler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	urlQueryParam := r.URL.Query().Get("url")
	if urlQueryParam == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	var err error
	if strings.Contains(urlPath, "/shorten") {
		err = h.repo.DeleteByShortUrl(r.Context(), urlQueryParam)
	} else {
		err = h.repo.DeleteByOriginalUrl(r.Context(), urlQueryParam)
	}

	// TODO: Validate error types of sqlc
	if err != nil {
		log.Printf("Failed to get URL: %v", err)
		http.Error(w, "Failed to get URL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
