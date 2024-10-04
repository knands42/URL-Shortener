package handler

import (
	"context"
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
// @Failure 404 {object} utils.NotFoundErrorResponse
// @Router /url [delete]
func (h *Handler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracing.Start(r.Context(), "DeleteURL")
	defer span.End()

	urlQueryParam := r.URL.Query().Get("url")
	urlTypeQuertParam := r.URL.Query().Get("type")
	if urlQueryParam == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	var err error
	if urlTypeQuertParam == URL_TYPE_SHORT {
		err = h.deleteUsingShortUrl(ctx, urlQueryParam)
	} else {
		err = h.deleteUsingOriginalUrl(ctx, urlQueryParam)
	}

	if err != nil {
		notFound(w, err, "URL not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deleteUsingShortUrl(ctx context.Context, shortUrl string) error {
	hash := h.extractHashFromUrl(shortUrl)
	err := h.repo.DeleteByHash(ctx, hash)
	if err != nil {
		return err
	}
	err = h.deleteFromCache(ctx, hash)
	if err != nil {
		return err
	}
	return h.deleteFromCache(ctx, hash+":metadata")
}

func (h *Handler) deleteUsingOriginalUrl(ctx context.Context, originalUrl string) error {
	data, err := h.repo.GetByOriginalUrl(ctx, originalUrl)
	if err != nil {
		return err
	}
	return h.deleteUsingShortUrl(ctx, data.Hash)
}
