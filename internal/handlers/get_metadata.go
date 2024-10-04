package handler

import (
	"context"
	"encoding/json"
	"knands42/url-shortener/internal/database/repo"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type GetMetadataResponse struct {
	OriginalUrl    string `json:"original_url"`
	ShortUrl       string `json:"short_url"`
	NumberOfAccess int32  `json:"number_of_access"`
}

func (h *Handler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracing.Start(r.Context(), "GetMetadata")
	defer span.End()

	urlTypeQueryParam := r.URL.Query().Get("type")
	urlPathParam := chi.URLParam(r, "url")

	var err error
	var resp GetMetadataResponse
	if urlTypeQueryParam == URL_TYPE_ORIGINAL {
		resp, err = h.getMetadataFromUrl(ctx, urlPathParam, false)
	} else {
		resp, err = h.getMetadataFromUrl(ctx, urlPathParam, true)
	}

	if err != nil {
		notFound(w, err, "URL not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) getMetadataFromUrl(ctx context.Context, url string, isShortUrl bool) (GetMetadataResponse, error) {
	ctx, span := h.tracing.Start(ctx, "GetMetadataFromShortenedUrl")
	defer span.End()

	var hash string
	if isShortUrl {
		hash = h.extractHashFromUrl(url)
	}

	cacheKey := buildCacheKeyForMedatada(url, hash, isShortUrl)
	cacheValue, err := h.getMetadataFromCache(ctx, cacheKey)
	if err != nil {
		log.Printf("Cache miss for %s: %v", hash, err.Error())
		resultDB, err := h.getMetadataFromRepo(ctx, url, hash, isShortUrl)
		if err != nil {
			return GetMetadataResponse{}, err
		}

		err = h.persistMetaddataIntoCache(ctx, cacheKey, resultDB)
		if err != nil {
			return GetMetadataResponse{}, err
		}

		return GetMetadataResponse{
			OriginalUrl:    resultDB.OriginalUrl,
			ShortUrl:       "https://me.li/" + resultDB.Hash,
			NumberOfAccess: resultDB.NumberOfAccess,
		}, nil
	}

	return GetMetadataResponse{
		OriginalUrl:    cacheValue.OriginalUrl,
		ShortUrl:       cacheValue.ShortUrl,
		NumberOfAccess: cacheValue.NumberOfAccess,
	}, nil
}

func buildCacheKeyForMedatada(url, hash string, isShort bool) string {
	if isShort {
		return hash + ":short"
	}
	return url + ":original"
}

func (h *Handler) getMetadataFromRepo(ctx context.Context, url, hash string, isShortUrl bool) (repo.ShortenedUrl, error) {
	if isShortUrl {
		return h.repo.GetByHash(ctx, hash)
	} else {
		return h.repo.GetByOriginalUrl(ctx, url)
	}
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
		OriginalUrl:    resultDB.OriginalUrl,
		ShortUrl:       "https://me.li/" + resultDB.Hash,
		NumberOfAccess: resultDB.NumberOfAccess,
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
