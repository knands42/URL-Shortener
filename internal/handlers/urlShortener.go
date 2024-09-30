package handler

import (
	"crypto/sha256"
	"encoding/json"
	"math/big"
	"net/http"
)

type GenerateShortURLRequest struct {
	Input  string `json:"input"`
	Length int    `json:"length"`
}

type GenerateShortURLResponse struct {
	ShortURL string `json:"short_url"`
}

func (h *Handler) GenerateShortURL(w http.ResponseWriter, r *http.Request) {
	var req GenerateShortURLRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// TODO: Add validation for the input and length fields

	input := req.Input
	length := req.Length

	if length < 6 || length > 12 {
		http.Error(w, "Length must be between 6 and 12", http.StatusBadRequest)
		return
	}

	// Hash the input using SHA256
	hash := sha256.Sum256([]byte(input))

	// Encode the hash to base62 instead of trimming the final result
	base62Hash := base62Encode(hash[:])

	if len(base62Hash) >= length {
		resp := GenerateShortURLResponse{
			ShortURL: base62Hash[:length],
		}

		json.NewEncoder(w).Encode(resp)
	} else {
		http.Error(w, "Length must be between 6 and 12", http.StatusBadRequest)
		return
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
