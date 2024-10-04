package handler

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"knands42/url-shortener/internal/database/repo"
	"knands42/url-shortener/internal/utils"
	"log"
	"net/http"

	_ "knands42/url-shortener/docs"

	"github.com/asaskevich/govalidator"
)

type GenerateShortURLRequest struct {
	Input string `json:"input" valid:"required,url" example:"https://www.google.com"`
}

type GenerateShortURLResponse struct {
	ShortURL string `json:"short_url" example:"https://me.li/abc123"`
}

// @Summary Generate a short URL
// @Description Generate a short URL from the input URL
// @Tags URL
// @Accept json
// @Produce json
// @Param input body GenerateShortURLRequest true "Input URL"
// @Success 201 {object} GenerateShortURLResponse
// @Failure 400 {object} utils.BadRequestErrorResponse
// @Failure 409 {object} utils.ConflictErrorResponse
// @Router /shorten [post]
func (h *Handler) GenerateShortURL(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracing.Start(r.Context(), "GenerateShortURL")
	defer span.End()

	var req GenerateShortURLRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	_, err = govalidator.ValidateStruct(req)
	if err != nil {
		validatorErrorResponse := utils.BadRequestErrorResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
		log.Printf("Validation error: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validatorErrorResponse)
		return
	}

	input := req.Input
	length := 7

	base62Hash := generateHash(input, length)

	_, err = h.repo.CreateHash(
		ctx,
		repo.CreateHashParams{
			OriginalUrl: input,
			Hash:        base62Hash[:length],
		},
	)
	if err != nil {
		conflict(w, err)
		return
	}

	hash := base62Hash[:length]
	h.saveIntoCache(ctx, hash, input)

	generateShortURLResponse := GenerateShortURLResponse{
		ShortURL: "https://me.li/" + hash,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(generateShortURLResponse)
}

func conflict(w http.ResponseWriter, err error) {
	errorResponse := utils.ConflictErrorResponse{
		Status:  http.StatusConflict,
		Message: "URL already exists",
	}
	log.Printf("URL already exists: %v", err.Error())
	w.WriteHeader(http.StatusConflict)
	json.NewEncoder(w).Encode(errorResponse)
}

func generateHash(input string, length int) string {

	// Hash the input using SHA256 to avoid collisions
	hash := sha256.Sum256([]byte(input))
	// Encode the hash to base62 instead of using raw hexadecimals limited to 16 characters
	base62Hash := base64.StdEncoding.EncodeToString(hash[:])

	return base62Hash[:length]
}
