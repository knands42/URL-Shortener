package handler

import (
	"context"
	"knands42/url-shortener/internal/database/repo"
	"log"
	"net/http"
)

// @Summary Get a URL entry
// @Description Get the original url from the shortened url and redirect to it
// @Tags URL
// @Accept json
// @Produce json
// @Param url query string true "URL to get metadata for"
// @Failure 404 {object} utils.NotFoundErrorResponse
// @Router /url [get]
func (h *Handler) GetOriginalUrl(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracing.Start(r.Context(), "GetUrl")
	defer span.End()

	urlQueryParam := r.URL.Query().Get("url")

	hash := h.extractHashFromUrl(urlQueryParam)
	var err error
	var resultDB repo.ShortenedUrl
	var cacheValue string
	var originalUrl string

	cacheValue, err = h.checkCacheFirst(ctx, hash)
	if err != nil {
		log.Printf("Cache miss for %s: %v", hash, err.Error())
		resultDB, err = h.getOriginalUrlFromRepo(ctx, hash)
		if err != nil {
			notFound(w, err, "URL not found")
			return
		}
		h.saveIntoCache(ctx, resultDB.Hash, resultDB.OriginalUrl)
		originalUrl = resultDB.OriginalUrl
	} else {
		originalUrl = cacheValue
	}

	w.WriteHeader(http.StatusOK)
	http.Redirect(w, r, originalUrl, http.StatusMovedPermanently)
}

func (h *Handler) getOriginalUrlFromRepo(ctx context.Context, valueParameter string) (repo.ShortenedUrl, error) {
	ctx, span := h.tracing.Start(ctx, "GetUrlFromRepo")
	defer span.End()

	return h.repo.GetByHash(ctx, valueParameter)
}
