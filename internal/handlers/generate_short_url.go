package handler

import (
	"crypto/sha256"
	"encoding/json"
	"knands42/url-shortener/internal/database/repo"
	"knands42/url-shortener/internal/utils"
	"log"
	"math/big"
	"net/http"

	_ "knands42/url-shortener/docs"

	"github.com/asaskevich/govalidator"
)

type GenerateShortURLRequest struct {
	Input string `json:"input" valid:"required,url"`
}

type GenerateShortURLResponse struct {
	ShortURL string `json:"short_url"`
}

// GenerateShortURL generates a short URL from the input URL
// @Summary Generate a short URL
func (h *Handler) GenerateShortURL(w http.ResponseWriter, r *http.Request) {
	var req GenerateShortURLRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	_, err = govalidator.ValidateStruct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := req.Input
	length := 6

	resp := generateFinalURL(input, length)

	_, err = h.repo.CreateShortUrl(
		r.Context(),
		repo.CreateShortUrlParams{
			OriginalUrl: input,
			ShortUrl:    resp.ShortURL,
		},
	)
	if err != nil {
		errorResponse := utils.ErrorResponse{
			Status:  http.StatusConflict,
			Message: "URL already exists",
		}
		log.Printf("URL already exists: %v", err.Error())
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func generateFinalURL(input string, length int) GenerateShortURLResponse {

	// Hash the input using SHA256 to avoid collisions
	hash := sha256.Sum256([]byte(input))
	// Encode the hash to base62 instead of using raw hexadecimals limited to 16 characters
	base62Hash := base62Encode(hash[:])

	return GenerateShortURLResponse{
		ShortURL: "https://me.li/" + base62Hash[:length],
	}
}

const base62Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

/*
Encodes a byte slice to a base62 string composed
of the following characters: [A-Za-z0-9].
*/
func base62Encode(input []byte) string {
	bigInt := new(big.Int).SetBytes(input)
	result := make([]byte, 0)
	base := big.NewInt(62)

	for bigInt.Cmp(big.NewInt(0)) > 0 {
		remainder := new(big.Int)
		bigInt.DivMod(bigInt, base, remainder)
		result = append(result, base62Chars[remainder.Int64()])
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}