package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
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
)
var dbConnection, _ = dbConfig.Connect(ctx)
var repository = repo.New(dbConnection)
var handlers = handler.NewHandler(repository)
var testServer = server.NewServer(r, handlers)

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

	b, _ := ioutil.ReadAll(resp.Body)
	_ = string(b)

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check the response body
	var responseBody map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if data, _ := responseBody["short_url"]; data != "https://me.li/BLvbbc" {
		t.Errorf("Expected response body to contain 'short_url' key")
	}
}
