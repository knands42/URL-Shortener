package handler

import (
	"encoding/json"
	"knands42/url-shortener/internal/utils"
	"log"
	"net/http"
)

// @Summary Delete a URL entry
// @Description Delete a URL entry by providing either the original URL or the short URL
// @Tags URL
// @Accept json
// @Produce json
// @Param url query string true "URL to be deleted"
// @Param type query string true "Type of URL to be deleted (short_url or original_url)"
// @Success 204
// @Failure 500 {object} utils.ErrorResponse
// @Router /url [delete]
func (h *Handler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracing.Start(r.Context(), "GenerateShortURL")
	defer span.End()

	urlQueryParam := r.URL.Query().Get("url")
	urlTypeQuertParam := r.URL.Query().Get("type")
	if urlQueryParam == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	var err error
	if urlTypeQuertParam == URL_TYPE_SHORT {
		hash := h.extractHashFromUrl(urlQueryParam)
		err = h.repo.DeleteByHash(ctx, hash)
	} else {
		err = h.repo.DeleteByOriginalUrl(ctx, urlQueryParam)
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
