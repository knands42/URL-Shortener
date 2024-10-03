package handler

import "context"

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

func (h *Handler) writeThroughCache(
	ctx context.Context,
	hash string,
	url string,
) {
	cacheShortKey := hash
	cacheOriginalKey := url
	h.saveIntoCache(ctx, cacheShortKey, url)
	h.saveIntoCache(ctx, cacheOriginalKey, hash)
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
