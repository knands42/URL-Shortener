package handler

import (
	"context"
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
// @Param url query string true "URL to be deleted" "https://www.google.com"
// @Param type query string true "Type of URL to be deleted (short_url or original_url)" "original_url"
// @Success 200 {object} GetURLResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /url [get]
func (h *Handler) GetUrl(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracing.Start(r.Context(), "GetUrl")
	defer span.End()

	urlQueryParam := r.URL.Query().Get("url")
	urlTypeQueryParam := r.URL.Query().Get("type")
	if urlQueryParam == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	var err error
	var resultDB repo.ShortenedUrl
	var cacheHit string
	var resp = GetURLResponse{}

	cacheKey, hash := h.getCacheKeyAndHash(ctx, urlTypeQueryParam, urlQueryParam)
	cacheHit, err = h.checkCacheFirst(ctx, cacheKey)
	if err != nil {
		log.Printf("Cache miss: %v", err.Error())
		resultDB, err = h.getFromRepo(ctx, urlTypeQueryParam, hash, urlQueryParam)
		if err != nil {
			notFound(w, err)
			return
		}
		h.saveIntoCache(ctx, cacheKey, h.getCacheValue(urlTypeQueryParam, resultDB))
		h.populateResultFromDB(&resp, resultDB)
	} else {
		h.populateResultFromCache(&resp, urlTypeQueryParam, cacheHit, urlQueryParam, hash)
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func notFound(w http.ResponseWriter, err error) {
	errorResponse := utils.ErrorResponse{
		Status:  http.StatusNotFound,
		Message: "URL not found",
	}
	log.Printf("URL not found: %v", err.Error())
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(errorResponse)
}

func (h *Handler) getFromRepo(ctx context.Context, urlTypeQueryParam, hash, urlQueryParam string) (repo.ShortenedUrl, error) {
	ctx, span := h.tracing.Start(ctx, "GetUrlFromRepo")
	defer span.End()

	if urlTypeQueryParam == URL_TYPE_SHORT {
		return h.repo.GetByHash(ctx, hash)
	}
	return h.repo.GetByOriginalUrl(ctx, urlQueryParam)
}

func (h *Handler) getCacheValue(urlTypeQueryParam string, resultDB repo.ShortenedUrl) string {
	if urlTypeQueryParam == URL_TYPE_SHORT {
		return resultDB.OriginalUrl
	}
	return resultDB.Hash
}

func (h *Handler) populateResultFromCache(resp *GetURLResponse, urlTypeQueryParam, cacheHit, urlQueryParam, hash string) {
	if urlTypeQueryParam == URL_TYPE_SHORT {
		resp.OriginalUrl = cacheHit
		resp.ShortUrl = "https://me.li/" + hash
	} else {
		resp.OriginalUrl = urlQueryParam
		resp.ShortUrl = cacheHit
	}
}

func (h *Handler) populateResultFromDB(resp *GetURLResponse, resultDB repo.ShortenedUrl) {
	resp.OriginalUrl = resultDB.OriginalUrl
	resp.ShortUrl = "https://me.li/" + resultDB.Hash
}
