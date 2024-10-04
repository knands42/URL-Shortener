package handler

import (
	"context"
	"knands42/url-shortener/internal/database/repo"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type GetURLResponse struct {
	OriginalUrl string `json:"original_url"`
	ShortUrl    string `json:"short_url"`
}

// @Summary Get a URL entry
// @Description Get the original url from the shortened url and redirect to it
// @Tags URL
// @Accept json
// @Produce json
// @Success 200 {object} GetURLResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /url/{url} [get]
func (h *Handler) GetOriginalUrl(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracing.Start(r.Context(), "GetUrl")
	defer span.End()

	urlPathParam := chi.URLParam(r, "url")

	hash := h.extractHashFromUrl(urlPathParam)
	var err error
	var resultDB repo.ShortenedUrl
	var cacheValue string
	var result = GetURLResponse{}

	cacheValue, err = h.checkCacheFirst(ctx, hash)
	if err != nil {
		log.Printf("Cache miss for %s: %v", hash, err.Error())
		resultDB, err = h.getOriginalUrlFromRepo(ctx, hash)
		if err != nil {
			notFound(w, err, "URL not found")
			return
		}
		h.saveIntoCache(ctx, resultDB.Hash, resultDB.OriginalUrl)
		populateResultFromDB(&result, resultDB)
	} else {
		populateResultFromCache(&result, hash, cacheValue)
	}

	w.WriteHeader(http.StatusOK)
	http.Redirect(w, r, result.OriginalUrl, http.StatusMovedPermanently)
}

func (h *Handler) getOriginalUrlFromRepo(ctx context.Context, valueParameter string) (repo.ShortenedUrl, error) {
	ctx, span := h.tracing.Start(ctx, "GetUrlFromRepo")
	defer span.End()

	return h.repo.GetByHash(ctx, valueParameter)
}

func populateResultFromCache(resp *GetURLResponse, key, value string) {
	resp.OriginalUrl = value
	resp.ShortUrl = "https://me.li/" + key
}

func populateResultFromDB(resp *GetURLResponse, resultDB repo.ShortenedUrl) {
	resp.OriginalUrl = resultDB.OriginalUrl
	resp.ShortUrl = "https://me.li/" + resultDB.Hash
}
