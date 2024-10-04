package handler

import (
	"context"
	"encoding/json"
	"knands42/url-shortener/internal/database/repo"
	"knands42/url-shortener/internal/utils"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
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
// @Param type query string true "Type of URL to be deleted (short_url or original_url)" "original_url"
// @Success 200 {object} GetURLResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /url/{url} [get]
func (h *Handler) GetUrl(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracing.Start(r.Context(), "GetUrl")
	defer span.End()

	urlTypeQueryParam := r.URL.Query().Get("type")
	urlPathParam := chi.URLParam(r, "url")

	queryParamToSearch := h.extractHashFromUrlIfThereIsAny(urlTypeQueryParam, urlPathParam)
	var err error
	var resultDB repo.ShortenedUrl
	var cacheValue string
	var result = GetURLResponse{}

	cacheValue, err = h.checkCacheFirst(ctx, queryParamToSearch)
	if err != nil {
		log.Printf("Cache miss for %s: %v", queryParamToSearch, err.Error())
		resultDB, err = h.getFromRepo(ctx, urlTypeQueryParam, queryParamToSearch)
		if err != nil {
			notFound(w, err)
			return
		}
		h.writeThroughCache(ctx, resultDB.Hash, resultDB.OriginalUrl)
		populateResultFromDB(&result, resultDB)
	} else {
		populateResultFromCache(&result, urlTypeQueryParam, cacheValue, queryParamToSearch)
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
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) extractHashFromUrlIfThereIsAny(urlTypeQueryParam string, urlParam string) string {
	if urlTypeQueryParam == URL_TYPE_SHORT {
		return h.extractHashFromUrl(urlParam)
	} else {
		return urlParam
	}
}

func (h *Handler) getFromRepo(ctx context.Context, urlTypeQueryParam, valueParameter string) (repo.ShortenedUrl, error) {
	ctx, span := h.tracing.Start(ctx, "GetUrlFromRepo")
	defer span.End()

	if urlTypeQueryParam == URL_TYPE_SHORT {
		return h.repo.GetByHash(ctx, valueParameter)
	}
	return h.repo.GetByOriginalUrl(ctx, valueParameter)
}

func populateResultFromCache(resp *GetURLResponse, urlTypeQueryParam, cacheValue, cacheKey string) {
	if urlTypeQueryParam == URL_TYPE_SHORT {
		resp.OriginalUrl = cacheValue
		resp.ShortUrl = "https://me.li/" + cacheKey
	} else {
		resp.OriginalUrl = cacheKey
		resp.ShortUrl = "https://me.li/" + cacheValue
	}
}

func populateResultFromDB(resp *GetURLResponse, resultDB repo.ShortenedUrl) {
	resp.OriginalUrl = resultDB.OriginalUrl
	resp.ShortUrl = "https://me.li/" + resultDB.Hash
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
