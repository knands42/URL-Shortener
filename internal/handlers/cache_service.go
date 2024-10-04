package handler

import (
	"context"
	"encoding/json"
)

type URLMetadataCacheData struct {
	OriginalUrl    string
	ShortUrl       string
	NumberOfAccess int32
}

func (h *URLMetadataCacheData) marshal() (string, error) {
	jsonData, err := json.Marshal(h)
	return string(jsonData), err
}

func (h *URLMetadataCacheData) unmarshal(data string) error {
	err := json.Unmarshal([]byte(data), h)
	return err
}

func (h *Handler) checkCacheFirst(ctx context.Context, key string) (string, error) {
	ctx, span := h.tracing.Start(ctx, "GetUrlFromCache")
	defer span.End()

	url, err := h.cache.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return url, nil
}

func (h *Handler) saveIntoCache(ctx context.Context, key string, value string) error {
	ctx, span := h.tracing.Start(ctx, "SaveIntoCache")
	defer span.End()

	err := h.cache.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) deleteFromCache(ctx context.Context, cacheKey string) error {
	ctx, span := h.tracing.Start(ctx, "DeleteFromCache")
	defer span.End()

	err := h.cache.Del(ctx, cacheKey).Err()
	if err != nil {
		return err
	}

	return nil
}
