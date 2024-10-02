package handler

import (
	"encoding/json"
	"knands42/url-shortener/internal/database/repo"
	"knands42/url-shortener/internal/utils"
	"log"
	"net/http"
)

type GetURLResponse struct {
	OriginalUrl string `json:"original_url"`
	ShortUrl    string `json:"short_url"`
}

// @Summary Get a URL entry
// @Description Get a URL entry by providing either the original URL or the short URL
// @Tags URL
// @Accept json
// @Produce json
// @Param url query string true "URL to be deleted"
// @Param type query string true "Type of URL to be deleted (short_url or original_url)"
// @Success 200 {object} GetURLResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /url [get]
func (h *Handler) GetUrl(w http.ResponseWriter, r *http.Request) {
	urlQueryParam := r.URL.Query().Get("url")
	urlTypeQuertParam := r.URL.Query().Get("type")
	if urlQueryParam == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	var err error
	var resultDB repo.ShortenedUrl

	if urlTypeQuertParam == URL_TYPE_SHORT {
		hash := h.extractHashFromUrl(urlQueryParam)
		resultDB, err = h.repo.GetByHash(r.Context(), hash)
	} else {
		resultDB, err = h.repo.GetByOriginalUrl(r.Context(), urlQueryParam)
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

	resp := GetURLResponse{
		OriginalUrl: resultDB.OriginalUrl,
		ShortUrl:    "https://me.li/" + resultDB.Hash,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
