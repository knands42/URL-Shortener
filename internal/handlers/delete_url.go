package handler

import (
	"encoding/json"
	"knands42/url-shortener/internal/utils"
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

	if err != nil {
		errorResponse := utils.ErrorResponse{
			Status:  http.StatusNotFound,
			Message: "URL not found",
		}
		log.Printf("URL not found: %v", err.Error())
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
		return

	}

	w.WriteHeader(http.StatusNoContent)
}
