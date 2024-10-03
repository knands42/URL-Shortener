package handler

import "context"

func (h *Handler) getShortUrlCacheKey(hash string) string {
	return URL_TYPE_SHORT + "-" + hash
}

func (h *Handler) getOriginalUrlCacheKey(original_url string) string {
	return URL_TYPE_ORIGINAL + "-" + original_url
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

func (h *Handler) getCacheKeyAndHash(ctx context.Context, urlTypeQueryParam, urlQueryParam string) (string, string) {
	_, span := h.tracing.Start(ctx, "GetCacheKeyAndHash")
	defer span.End()

	if urlTypeQueryParam == URL_TYPE_SHORT {
		hash := h.extractHashFromUrl(urlQueryParam)
		return h.getShortUrlCacheKey(hash), hash
	}
	return h.getOriginalUrlCacheKey(urlQueryParam), ""
}
