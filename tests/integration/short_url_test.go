package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"knands42/url-shortener/internal/database"
	"knands42/url-shortener/internal/database/repo"
	handler "knands42/url-shortener/internal/handlers"
	"knands42/url-shortener/internal/server"
	"knands42/url-shortener/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

var r = chi.NewRouter()
var ctx = context.Background()
var config = utils.NewConfig("test")
var dbConfig = database.NewDBConfig(
	config.DBHost,
	config.DBUser,
	config.DBPassword,
	config.DBName,
	config.DBPort,
	config.SSLMode,
	config.TimeZone,
	config.MinConns,
	config.MaxConns,
)
var dbConnection, _ = dbConfig.Connect(ctx)
var repository = repo.New(dbConnection)
var handlers = handler.NewHandler(repository)
var testServer = server.NewServer(r, handlers)

func Test_delete_entry_by_short_url(t *testing.T) {
	ts := httptest.NewServer(testServer.Router)

	// Send a DELETE request to the /api/v1/shorten endpoint
	req, err := http.NewRequest("DELETE", ts.URL+"/api/v1/url?url=https://google.com&type=original_url", nil)
	if err != nil {
		t.Fatalf("Failed to send DELETE request: %v", err)
	}

	// Check the response status code
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send DELETE request: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, resp.StatusCode)
	}
}

func Test_create_short_url(t *testing.T) {
	ts := httptest.NewServer(testServer.Router)

	// Define the request payload
	payload := map[string]string{"input": "https://google.com"}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	// Send a POST request to the /api/v1/shorten endpoint
	resp, err := http.Post(ts.URL+"/api/v1/shorten", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check the response body
	var responseBody map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if data := responseBody["short_url"]; data != "https://me.li/BQRvJsg" {
		t.Errorf("Expected response body to contain 'short_url' key")
	}
}

func Test_get_entry_by_short_url(t *testing.T) {
	ts := httptest.NewServer(testServer.Router)

	// Send a GET request to the /api/v1/shorten endpoint
	req, err := http.NewRequest("GET", ts.URL+"/api/v1/url?url=https://me.li/BQRvJsg&type=short_url", nil)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}

	// Check the response status code
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, resp.StatusCode)
	}

	// Check the response body
	var responseBody map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if data := responseBody["original_url"]; data != "https://google.com" {
		t.Errorf("Expected response body to contain 'original_url' key")
	}

	if data := responseBody["short_url"]; data != "https://me.li/BQRvJsg" {
		t.Errorf("Expected response body to contain 'short_url' key")
	}
}

func Test_get_entry_by_original_url(t *testing.T) {
	ts := httptest.NewServer(testServer.Router)

	// Send a GET request to the /api/v1/shorten endpoint
	req, err := http.NewRequest("GET", ts.URL+"/api/v1/url?url=https://google.com&type=original_url", nil)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}

	// Check the response status code
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, resp.StatusCode)
	}

	// Check the response body
	var responseBody map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if data := responseBody["original_url"]; data != "https://google.com" {
		t.Errorf("Expected response body to contain 'original_url' key")
	}

	if data := responseBody["short_url"]; data != "https://me.li/BQRvJsg" {
		t.Errorf("Expected response body to contain 'short_url' key")
	}
}

func Test_delete_entry_by_original_url(t *testing.T) {
	ts := httptest.NewServer(testServer.Router)

	// Send a DELETE request to the /api/v1/shorten endpoint
	req, err := http.NewRequest("DELETE", ts.URL+"/api/v1/url?url=https://me.li/BQRvJsg&type=short_url", nil)
	if err != nil {
		t.Fatalf("Failed to send DELETE request: %v", err)
	}

	// Check the response status code
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send DELETE request: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, resp.StatusCode)
	}
}
