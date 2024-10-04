package handler

import (
	"encoding/json"
	"knands42/url-shortener/internal/database/repo"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type GetMetadataResponse struct {
	NumberOfAccess int32 `json:"number_of_access"`
}

func (h *Handler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracing.Start(r.Context(), "GetMetadata")
	defer span.End()

	urlTypeQueryParam := r.URL.Query().Get("type")
	urlPathParam := chi.URLParam(r, "url")

	queryParamToSearch := h.extractHashFromUrlIfThereIsAny(urlTypeQueryParam, urlPathParam)

	var err error
	var resultDB repo.ShortenedUrl
	var result = GetMetadataResponse{}

	resultDB, err = h.getFromRepo(ctx, urlTypeQueryParam, queryParamToSearch)

	if err != nil {
		notFound(w, err)
		return
	}

	result.NumberOfAccess = resultDB.NumberOfAccess

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
