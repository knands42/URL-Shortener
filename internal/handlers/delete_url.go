package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) DeleteURL(w http.ResponseWriter, r *http.Request) {
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

	var err error
	if urlType == URL_TYPE_SHORT {
		err = h.repo.DeleteByShortUrl(r.Context(), url)
	} else {
		err = h.repo.DeleteByOriginalUrl(r.Context(), url)
	}

	// TODO: Validate error types of sqlc
	if err != nil {
		log.Printf("Failed to get URL: %v", err)
		http.Error(w, "Failed to get URL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
