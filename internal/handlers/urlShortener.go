package handler

import (
	"crypto/sha256"
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/asaskevich/govalidator"
)

type GenerateShortURLRequest struct {
	Input  string `json:"input" valid:"required,url"`
	Length int    `json:"length" valid:"optional,numeric,range(6|12)"`
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

	_, err = govalidator.ValidateStruct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
		// TODO: create a custom error response
	}

	input := req.Input
	length := 10
	if req.Length != 0 {
		length = req.Length
	}

	resp := generateFinalHash(input, length)
	json.NewEncoder(w).Encode(resp)
}

func generateFinalHash(input string, length int) GenerateShortURLResponse {

	// Hash the input using SHA256
	hash := sha256.Sum256([]byte(input))

	// Encode the hash to base62 instead of trimming the final result
	base62Hash := base62Encode(hash[:])

	return GenerateShortURLResponse{
		ShortURL: base62Hash[:length],
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
