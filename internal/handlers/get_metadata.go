package handler

import (
	"context"
	"encoding/json"
	"knands42/url-shortener/internal/database/repo"
	"log"
	"net/http"
)

type GetMetadataResponse struct {
	OriginalUrl      string `json:"original_url" example:"https://www.google.com"`
	ShortUrl         string `json:"short_url" example:"https://me.li/abc123"`
	NumberOfAccesses int32  `json:"number_of_access" example:"0"`
	CreatedAt        string `json:"created_at" example:"2021-09-01 12:00:00"`
	UpdatedAt        string `json:"updated_at" example:"2021-09-01 12:00:00"`
}

// @Summary Get a URL entry
// @Description Get information about any URL entry by providing the original URL or the shortened URL
// @Tags URL
// @Accept json
// @Produce json
// @Param url query string true "URL to get metadata for"
// @Param type query string false "Type of URL (short_url or original_url)"
// @Success 200 {object} GetMetadataResponse
// @Failure 404 {object} utils.NotFoundErrorResponse
// @Router /url/metadata [get]
func (h *Handler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracing.Start(r.Context(), "GetMetadata")
	defer span.End()

	urlTypeQueryParam := r.URL.Query().Get("type")
	urlQueryParam := r.URL.Query().Get("url")

	var err error
	var resp GetMetadataResponse
	if urlTypeQueryParam == URL_TYPE_ORIGINAL {
		resp, err = h.getMetadataFromOriginalUrl(ctx, urlQueryParam)
	} else {
		resp, err = h.getMetadataFromShortUrl(ctx, urlQueryParam)
	}

	if err != nil {
		notFound(w, err, "URL not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) getMetadataFromShortUrl(ctx context.Context, url string) (GetMetadataResponse, error) {
	ctx, span := h.tracing.Start(ctx, "GetMetadataFromShortenedUrl")
	defer span.End()

	hash := h.extractHashFromUrl(url)

	cacheKey := hash + ":metadata"
	cacheValue, err := h.getMetadataFromCache(ctx, cacheKey)
	if err != nil {
		log.Printf("Cache miss for %s: %v", hash, err.Error())
		resultDB, err := h.repo.GetByHash(ctx, hash)
		if err != nil {
			return GetMetadataResponse{}, err
		}

		err = h.persistMetaddataIntoCache(ctx, cacheKey, resultDB)
		if err != nil {
			return GetMetadataResponse{}, err
		}

		return GetMetadataResponse{
			OriginalUrl:      resultDB.OriginalUrl,
			ShortUrl:         "https://me.li/" + resultDB.Hash,
			NumberOfAccesses: resultDB.NumberOfAccesses,
			CreatedAt:        resultDB.CreatedAt.Time.GoString(),
			UpdatedAt:        resultDB.UpdatedAt.Time.GoString(),
		}, nil
	}

	return GetMetadataResponse(cacheValue), nil
}

func (h *Handler) getMetadataFromOriginalUrl(ctx context.Context, url string) (GetMetadataResponse, error) {
	ctx, span := h.tracing.Start(ctx, "GetMetadataFromOriginalUrl")
	defer span.End()

	resultDB, err := h.repo.GetByOriginalUrl(ctx, url)
	if err != nil {
		return GetMetadataResponse{}, err
	}

	return GetMetadataResponse{
		OriginalUrl:      resultDB.OriginalUrl,
		ShortUrl:         "https://me.li/" + resultDB.Hash,
		NumberOfAccesses: resultDB.NumberOfAccesses,
		CreatedAt:        resultDB.CreatedAt.Time.GoString(),
		UpdatedAt:        resultDB.UpdatedAt.Time.GoString(),
	}, nil
}

func (h *Handler) getMetadataFromCache(ctx context.Context, key string) (URLMetadataCacheData, error) {
	url, err := h.checkCacheFirst(ctx, key)
	if err != nil {
		return URLMetadataCacheData{}, err
	}

	var urlMetadataCacheData URLMetadataCacheData
	err = urlMetadataCacheData.unmarshal(url)
	if err != nil {
		return URLMetadataCacheData{}, err
	}

	return urlMetadataCacheData, nil
}

func (h *Handler) persistMetaddataIntoCache(ctx context.Context, key string, resultDB repo.ShortenedUrl) error {
	ctx, span := h.tracing.Start(ctx, "PersistIntoCache")
	defer span.End()

	urlMetadataCacheData := URLMetadataCacheData{
		OriginalUrl:      resultDB.OriginalUrl,
		ShortUrl:         "https://me.li/" + resultDB.Hash,
		NumberOfAccesses: resultDB.NumberOfAccesses,
		CreatedAt:        resultDB.CreatedAt.Time.GoString(),
		UpdatedAt:        resultDB.UpdatedAt.Time.GoString(),
	}
	cacheData, err := urlMetadataCacheData.marshal()
	if err != nil {
		return err
	}

	err = h.saveIntoCache(ctx, key, cacheData)
	if err != nil {
		return err
	}

	return nil
}
